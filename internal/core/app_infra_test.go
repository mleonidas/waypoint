package core

import (
	"context"
	"testing"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	componentmocks "github.com/hashicorp/waypoint-plugin-sdk/component/mocks"
	"github.com/hashicorp/waypoint/internal/config"
	"github.com/stretchr/testify/require"
)

func TestAppInfra_happy(t *testing.T) {
	require := require.New(t)

	// Make our factory for platforms
	mock := &componentmocks.Infra{}
	factory := TestFactory(t, component.InfraType)
	TestFactoryRegister(t, factory, "test", mock)

	// Make our app
	app := TestApp(t, TestProject(t,
		WithConfig(config.TestConfig(t, testInfraConfig)),
		WithFactory(component.InfraType, factory),
		WithJobInfo(&component.JobInfo{Id: "hello"}),
	), "test")

	// Setup our value
	artifact := &componentmocks.Artifact{}
	artifact.On("Labels").Return(map[string]string{"foo": "foo"})
	mock.On("InfraFunc").Return(func() component.Artifact {
		return artifact
	})

	{
		// Destroy
		infra, err := app.Infra(context.Background())
		require.NoError(err)

		// Verify that we set the status properly
		require.Equal("foo", infra.Labels["foo"])
		require.Contains(infra.Labels, "waypoint/workspace")

		// Verify we have the ID set
		require.Equal("hello", infra.JobId)
	}
}

const testInfraConfig = `

project = "test"

app "test" {
    infra {
	    labels = { "foo" = "bar" }
        use "test" {
        }
    }

	build {
		labels = { "foo" = "bar" }
		use "test" {
			foo = labels["foo"]
		}
	}

	deploy {
		use "test" {}
	}
}
`
