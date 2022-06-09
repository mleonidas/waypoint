package singleprocess

import (
	"context"

	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) UpsertInfrastructure(
	ctx context.Context,
	req *pb.UpsertInfraRequest,
) (*pb.UpsertInfraResponse, error) {
	if err := ptypes.ValidateUpsertInfraRequest(req); err != nil {
		return nil, err
	}

	result := req.Infrastructure

	// If we have no ID, then we're inserting and need to generate an ID.
	insert := result.Id == ""
	if insert {
		// Get the next id
		id, err := server.Id()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "uuid generation failed: %s", err)
		}

		// Specify the id
		result.Id = id
	}

	if err := s.state(ctx).InfraPut(!insert, result); err != nil {
		return nil, err
	}

	return &pb.UpsertInfraResponse{Infrastructure: result}, nil
}
