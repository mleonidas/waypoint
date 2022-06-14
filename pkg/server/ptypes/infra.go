package ptypes

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"
)

// TestBuild returns a valid user for tests.
func TestInfra(t testing.T, src *pb.Infra) *pb.Infra {
	t.Helper()

	if src == nil {
		src = &pb.Infra{}
	}

	require.NoError(t, mergo.Merge(src, &pb.Infra{
		Id: "test",
	}))

	return src
}

// ValidateGetInfraRequest
func ValidateGetInfraRequest(v *pb.GetInfraRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Ref, validation.Required),
		validationext.StructField(&v.Ref, func() []*validation.FieldRules {
			return ValidateRefOperationRules(v.Ref)
		}),
	))
}

// ValidateGetLatestBuildRequest
func ValidateGetLatestInfraRequest(v *pb.GetLatestInfraRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Application, validation.Required),
	))
}

// ValidateUpsertBuildRequest
func ValidateUpsertInfraRequest(v *pb.UpsertInfraRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Infra, validation.Required),
		validationext.StructField(&v.Infra, func() []*validation.FieldRules {
			return ValidateInfraRules(v.Infra)
		}),
	))
}

// ValidateBuildRules
func ValidateInfraRules(v *pb.Infra) []*validation.FieldRules {
	return []*validation.FieldRules{
		validationext.StructField(&v.Application, func() []*validation.FieldRules {
			return []*validation.FieldRules{
				validation.Field(&v.Application.Application, validation.Required),
				validation.Field(&v.Application.Project, validation.Required),
			}
		}),

		validationext.StructField(&v.Workspace, func() []*validation.FieldRules {
			return []*validation.FieldRules{
				validation.Field(&v.Workspace.Workspace, validation.Required),
			}
		}),
	}
}

// ValidateListBuildsRequest
func ValidateListInfraRequest(v *pb.ListInfraRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validationext.StructField(&v.Application, func() []*validation.FieldRules {
			return []*validation.FieldRules{
				validation.Field(&v.Application.Application, validation.Required),
				validation.Field(&v.Application.Project, validation.Required),
			}
		})))
}
