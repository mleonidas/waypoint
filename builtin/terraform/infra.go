package terraform

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Infra struct {
	config InfraConfig
}

type InfraConfig struct {
	ModuleSrc string `hcl:"source"`
	VarFile   string `hcl:"var_path"`
}

func (i *Infra) InfraFunc() interface{} {
	return i.Infrastructure
}

func (i *Infra) Infrastructure(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	ui terminal.UI,
) (*Infrastructure, error) {
	var result Infrastructure

	sg := ui.StepGroup()
	defer sg.Wait()

	log.Debug("made it into infra plugin")
	s := sg.Add("Initializing terraform infra...")
	defer s.Abort()

	ulid, err := component.Id()

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate ULID: %s", err)
	}

	infraId := fmt.Sprintf("%s-%s", src.App, ulid)
	log.Debug(infraId)

	if i.config.ModuleSrc == "" {
		log.Debug("no module source configured")
	}
	result.Cluster = "test-cluster"
	result.ClusterId = "eks:something:somethign::"
	s.Done()

	return &result, nil
}

// Config implements Configurable
func (b *Infra) Config() (interface{}, error) {
	return &b.config, nil
}

func (i *Infra) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&InfraConfig{}), docs.FromFunc(i.InfraFunc()))

	if err != nil {
		return nil, err
	}

	doc.Description("Infrastructure provisioning with terraform")
	doc.Example(`
    infra {
        use "terraform" {
            source = "https://github.com/someorg/somemodule"
            var_file = "vars/prod.tfvars"
        }
    }
    `)

	return doc, nil
}

var (
	_ component.Infra        = (*Infra)(nil)
	_ component.Configurable = (*Infra)(nil)
	// _ component.Destroyer    = (*Infra)(nil)
	_ component.Documented = (*Infra)(nil)
)
