package cli

import (
	"context"

	"github.com/apex/log"
	"github.com/posener/complete"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

type InfraCommand struct {
	*baseCommand

	flagPush bool
}

func (c *InfraCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithMultiAppTargets(),
	); err != nil {
		return 1
	}

	err := c.DoApp(c.Ctx, func(ctx context.Context, app *clientpkg.App) error {
		app.UI.Output("Building infra %s...", app.Ref().Application, terminal.WithHeaderStyle())
		infraResult, err := app.Infra(ctx, &pb.Job_InfraOp{
			DisablePush: !c.flagPush,
		})
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}
		// Show input variable values used in infra
		app.UI.Output("Variables used:", terminal.WithHeaderStyle())
		log.Info("JobID: " + infraResult.Infra.JobId)
		resp, err := c.project.Client().GetJob(ctx, &pb.GetJobRequest{
			JobId: infraResult.Infra.JobId,
		})
		if err != nil {
			app.UI.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return ErrSentinel
		}
		tbl := fmtVariablesOutput(resp.VariableFinalValues)
		c.ui.Table(tbl)

		return nil
	})
	if err != nil {
		return 1
	}

	return 0
}

func (c *InfraCommand) Flags() *flag.Sets {
	return c.flagSet(flagSetOperation, func(set *flag.Sets) {
		f := set.NewSet("Command Options")
		f.BoolVar(&flag.BoolVar{
			Name:    "push",
			Target:  &c.flagPush,
			Default: true,
			Usage:   "TODO: fix this",
		})
	})
}

func (c *InfraCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *InfraCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *InfraCommand) Synopsis() string {
	return "Build a new versioned artifact from source"
}

func (c *InfraCommand) Help() string {
	return formatHelp(`
Usage: waypoint infra [options]
Alias: waypoint infra [options]

  Build a infra <somethign>

` + c.Flags().Help())
}
