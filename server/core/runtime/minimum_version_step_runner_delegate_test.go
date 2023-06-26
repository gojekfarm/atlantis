package runtime

import (
	"testing"

	"github.com/hashicorp/go-version"
	. "github.com/petergtz/pegomock"
	"github.com/gojekfarm/atlantis/server/core/runtime/mocks"
	"github.com/gojekfarm/atlantis/server/events/command"
	. "github.com/gojekfarm/atlantis/testing"
)

func TestRunMinimumVersionDelegate(t *testing.T) {
	RegisterMockTestingT(t)

	mockDelegate := mocks.NewMockRunner()

	tfVersion12, _ := version.NewVersion("0.12.0")
	tfVersion11, _ := version.NewVersion("0.11.15")

	// these stay the same for all tests
	extraArgs := []string{"extra", "args"}
	envs := map[string]string{}
	path := ""

	expectedOut := "some valid output from delegate"

	t.Run("default version success", func(t *testing.T) {
		subject := &MinimumVersionStepRunnerDelegate{
			defaultTfVersion: tfVersion12,
			minimumVersion:   tfVersion12,
			delegate:         mockDelegate,
		}

		ctx := command.ProjectContext{}

		When(mockDelegate.Run(ctx, extraArgs, path, envs)).ThenReturn(expectedOut, nil)

		output, err := subject.Run(
			ctx,
			extraArgs,
			path,
			envs,
		)

		Equals(t, expectedOut, output)
		Ok(t, err)
	})

	t.Run("ctx version success", func(t *testing.T) {
		subject := &MinimumVersionStepRunnerDelegate{
			defaultTfVersion: tfVersion11,
			minimumVersion:   tfVersion12,
			delegate:         mockDelegate,
		}

		ctx := command.ProjectContext{
			TerraformVersion: tfVersion12,
		}

		When(mockDelegate.Run(ctx, extraArgs, path, envs)).ThenReturn(expectedOut, nil)

		output, err := subject.Run(
			ctx,
			extraArgs,
			path,
			envs,
		)

		Equals(t, expectedOut, output)
		Ok(t, err)
	})

	t.Run("default version failure", func(t *testing.T) {
		subject := &MinimumVersionStepRunnerDelegate{
			defaultTfVersion: tfVersion11,
			minimumVersion:   tfVersion12,
			delegate:         mockDelegate,
		}

		ctx := command.ProjectContext{}

		output, err := subject.Run(
			ctx,
			extraArgs,
			path,
			envs,
		)

		mockDelegate.VerifyWasCalled(Never())

		Equals(t, "Version: 0.11.15 is unsupported for this step. Minimum version is: 0.12.0", output)
		Ok(t, err)
	})

	t.Run("ctx version failure", func(t *testing.T) {
		subject := &MinimumVersionStepRunnerDelegate{
			defaultTfVersion: tfVersion12,
			minimumVersion:   tfVersion12,
			delegate:         mockDelegate,
		}

		ctx := command.ProjectContext{
			TerraformVersion: tfVersion11,
		}

		output, err := subject.Run(
			ctx,
			extraArgs,
			path,
			envs,
		)

		mockDelegate.VerifyWasCalled(Never())

		Equals(t, "Version: 0.11.15 is unsupported for this step. Minimum version is: 0.12.0", output)
		Ok(t, err)
	})

}
