package handlertest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func init() {
	tests["infra"] = []testFunc{
		TestServiceInfra,
	}
}
func TestServiceInfra(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	_, client := factory(t)

	// Simplify writing tests
	type Req = pb.UpsertInfraRequest

	t.Run("create and update", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertInfra(ctx, &Req{
			Infra: serverptypes.TestValidInfra(t, nil),
		})
		require.NoError(err)
		require.NotNil(resp)
		result := resp.Infra
		require.NotEmpty(result.Id)

		// Let's write some data
		result.Status = server.NewStatus(pb.Status_RUNNING)
		resp, err = client.UpsertInfra(ctx, &Req{
			Infra: result,
		})
		require.NoError(err)
		require.NotNil(resp)
		result = resp.Infra
		require.NotNil(result.Status)
		require.Equal(pb.Status_RUNNING, result.Status.State)
	})

	t.Run("update non-existent", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertInfra(ctx, &Req{
			Infra: serverptypes.TestValidInfra(t, &pb.Infra{Id: "nope"}),
		})
		require.Error(err)
		require.Nil(resp)
		st, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, st.Code())
	})
}
