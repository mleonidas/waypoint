package core

import (
	"context"
	"errors"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/opaqueany"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/plugin"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// Infra infra the artifact from source for this app.
// TODO(mitchellh): test
func (a *App) Infra(ctx context.Context, optFuncs ...InfraOption) (
	*pb.Infra,
	error,
) {
	opts, err := newInfraOptions(optFuncs...)
	if err != nil {
		return nil, err
	}

	// Render the config
	c, err := componentCreatorMap[component.InfraType].Create(ctx, a, nil)
	if err != nil {
		return nil, err
	}

	defer c.Close()

	cr, err := componentCreatorMap[component.RegistryType].Create(ctx, a, nil)
	if err != nil {
		if status.Code(err) == codes.Unimplemented {
			cr = nil
			err = nil
		} else {
			return nil, err
		}
	}

	if cr != nil {
		defer cr.Close()
	}

	// First we do the infra
	_, msg, err := a.doOperation(ctx, a.logger.Named("infra"), &infraOperation{
		Component:   c,
		Registry:    cr,
		HasRegistry: cr != nil,
	})

	if err != nil {
		return nil, err
	}
	infra := msg.(*pb.Infra)

	// If we're not pushing, then we're done!
	if !opts.Push {
		return infra, nil
	}

	// We're also pushing to a registry, so invoke that.
	return infra, err
}

// Name returns the name of the operation
func (op *infraOperation) Name() string {
	return "infra"
}

// InfraOption is used to configure a Infra
type InfraOption func(*infraOptions) error

// InfraWithPush sets whether or not the infra will push. The default
// is for the infra to push.
func InfraWithPush(v bool) InfraOption {
	return func(opts *infraOptions) error {
		opts.Push = v
		return nil
	}
}

type infraOptions struct {
	Push bool
}

func defaultInfraOptions() *infraOptions {
	return &infraOptions{
		Push: true,
	}
}

func newInfraOptions(opts ...InfraOption) (*infraOptions, error) {
	def := defaultInfraOptions()
	for _, f := range opts {
		if err := f(def); err != nil {
			return nil, err
		}
	}

	return def, def.Validate()
}

func (opts *infraOptions) Validate() error {
	return nil
}

// infra implements the operation interface.
type infraOperation struct {
	Component *Component
	Registry  *Component
	Infra     *pb.Infra

	HasRegistry bool
}

func (op *infraOperation) Init(app *App) (proto.Message, error) {
	return &pb.Infra{
		Application: app.ref,
		Workspace:   app.workspace,
		Component:   op.Component.Info,
	}, nil
}

func (op *infraOperation) Hooks(app *App) map[string][]*config.Hook {
	return op.Component.hooks
}

func (op *infraOperation) Labels(app *App) map[string]string {
	return op.Component.labels
}

func (op *infraOperation) Upsert(
	ctx context.Context,
	client pb.WaypointClient,
	msg proto.Message,
) (proto.Message, error) {
	resp, err := client.UpsertInfra(ctx, &pb.UpsertInfraRequest{
		Infra: msg.(*pb.Infra),
	})

	if err != nil {
		return nil, err
	}

	return resp.Infra, nil
}

func (op *infraOperation) Do(ctx context.Context, log hclog.Logger, app *App, _ proto.Message) (interface{}, error) {
	args := []argmapper.Arg{
		argmapper.Named("HasRegistry", op.HasRegistry),
	}

	// If there is a registry defined and it implements RegistryAccess...
	if op.Registry != nil {
		if ra, ok := op.Registry.Value.(component.RegistryAccess); ok && ra.AccessInfoFunc() != nil {
			raw, err := app.callDynamicFunc(ctx, log, nil, op.Component, ra.AccessInfoFunc())
			if err == nil {
				args = append(args, argmapper.Typed(raw))

				if pm, ok := raw.(interface {
					TypedAny() *opaqueany.Any
				}); ok {
					any := pm.TypedAny()

					// ... which we make available to infra plugin.
					args = append(args, plugin.ArgNamedAny("access_info", any))
					log.Debug("injected access info")
				} else {
					log.Error("unexpected response type from callDynamicFunc", "type", hclog.Fmt("%T", raw))
					return nil, errors.New("AccessInfoFunc didn't provide a typed any")
				}
			} else {
				log.Error("error calling dynamic func", "error", err)
				return nil, err
			}
		} else {
			if ok && ra != nil && ra.AccessInfoFunc() == nil {
				return nil, status.Error(codes.Internal, "The plugin requested does not "+
					"define an AccessInfoFunc() in its Registry plugin. This is an internal "+
					"error and should be reported to the author of the plugin.")
			}
		}
	}

	return app.callDynamicFunc(ctx,
		log,
		(*component.Artifact)(nil),
		op.Component,
		op.Component.Value.(component.Infra).InfraFunc(),
		args...,
	)
}

func (op *infraOperation) StatusPtr(msg proto.Message) **pb.Status {
	return &(msg.(*pb.Infra).Status)
}

func (op *infraOperation) ValuePtr(msg proto.Message) (**opaqueany.Any, *string) {
	v := msg.(*pb.Infra)
	if v.Artifact == nil {
		v.Artifact = &pb.Artifact{}
	}

	return &v.Artifact.Artifact, &v.Artifact.ArtifactJson
}

var _ operation = (*infraOperation)(nil)
