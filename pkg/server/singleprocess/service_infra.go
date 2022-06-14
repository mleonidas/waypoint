package singleprocess

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

func (s *Service) UpsertInfra(
	ctx context.Context,
	req *pb.UpsertInfraRequest,
) (*pb.UpsertInfraResponse, error) {
	if err := serverptypes.ValidateUpsertInfraRequest(req); err != nil {
		return nil, err
	}

	result := req.Infra

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

	return &pb.UpsertInfraResponse{Infra: result}, nil
}

// GetDeployment returns a Deployment based on ID
func (s *Service) GetInfra(
	ctx context.Context,
	req *pb.GetInfraRequest,
) (*pb.Infra, error) {
	if err := serverptypes.ValidateGetInfraRequest(req); err != nil {
		return nil, err
	}

	d, err := s.state(ctx).InfraGet(req.Ref)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (s *Service) GetLatestInfra(
	ctx context.Context,
	req *pb.GetLatestInfraRequest,
) (*pb.Infra, error) {
	if err := serverptypes.ValidateGetLatestInfraRequest(req); err != nil {
		return nil, err
	}

	return s.state(ctx).InfraLatest(req.Application, req.Workspace)
}

func (s *Service) ListInfras(
	ctx context.Context,
	req *pb.ListInfraRequest,
) (*pb.ListInfraResponse, error) {
	if err := serverptypes.ValidateListInfraRequest(req); err != nil {
		return nil, err
	}

	result, err := s.state(ctx).InfraList(req.Application,
		serverstate.ListWithWorkspace(req.Workspace),
		serverstate.ListWithOrder(req.Order),
	)
	if err != nil {
		return nil, err
	}

	return &pb.ListInfraResponse{Infra: result}, nil
}
