package runner

import (
	"context"

	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func (r *Runner) executeInfraOp(
	ctx context.Context,
	job *pb.Job,
	project *core.Project,
) (*pb.Job_Result, error) {
	app, err := project.App(job.Application.Application)
	if err != nil {
		return nil, err
	}

	op, ok := job.Operation.(*pb.Job_Infra)
	if !ok {
		// this shouldn't happen since the call to this function is gated
		// on the above type match.
		panic("operation not expected type")
	}

	infra, err := app.Infra(ctx, core.InfraWithPush(!op.Infra.DisablePush))
	if err != nil {
		return nil, err
	}

	return &pb.Job_Result{
		Infra: &pb.Job_InfraResult{
			Infra: infra,
		},
	}, nil
}
