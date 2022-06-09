package boltdbstate

import (
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

var infraOp = &appOperation{
	Struct: (*pb.Infra)(nil),
	Bucket: []byte("infra"),
}

func init() {
	infraOp.register()
}

// InfraPut inserts or updates a infra record.
func (s *State) InfraPut(update bool, b *pb.Infra) error {
	return infraOp.Put(s, update, b)
}

// InfraGet gets a infra by ref.
func (s *State) InfraGet(ref *pb.Ref_Operation) (*pb.Infra, error) {
	result, err := infraOp.Get(s, ref)
	if err != nil {
		return nil, err
	}

	return result.(*pb.Infra), nil
}

func (s *State) InfraList(
	ref *pb.Ref_Application,
	opts ...serverstate.ListOperationOption,
) ([]*pb.Infra, error) {
	raw, err := infraOp.List(s, serverstate.BuildListOperationOptions(ref, opts...))
	if err != nil {
		return nil, err
	}

	result := make([]*pb.Infra, len(raw))
	for i, v := range raw {
		result[i] = v.(*pb.Infra)
	}

	return result, nil
}

// InfraLatest gets the latest infra that was completed successfully.
func (s *State) InfraLatest(
	ref *pb.Ref_Application,
	ws *pb.Ref_Workspace,
) (*pb.Infra, error) {
	result, err := infraOp.Latest(s, ref, ws)
	if result == nil || err != nil {
		return nil, err
	}

	return result.(*pb.Infra), nil
}
