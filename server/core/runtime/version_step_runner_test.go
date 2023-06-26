package runtime

import (
	"testing"

	"github.com/hashicorp/go-version"
	. "github.com/petergtz/pegomock"
	"github.com/gojekfarm/atlantis/server/core/terraform/mocks"
	"github.com/gojekfarm/atlantis/server/events/command"
	"github.com/gojekfarm/atlantis/server/events/models"
	"github.com/gojekfarm/atlantis/server/logging"
	. "github.com/gojekfarm/atlantis/testing"
)

func TestRunVersionStep(t *testing.T) {
	RegisterMockTestingT(t)
	logger := logging.NewNoopLogger(t)
	workspace := "default"

	context := command.ProjectContext{
		Log:                logger,
		EscapedCommentArgs: []string{"comment", "args"},
		Workspace:          workspace,
		RepoRelDir:         ".",
		User:               models.User{Username: "username"},
		Pull: models.PullRequest{
			Num: 2,
		},
		BaseRepo: models.Repo{
			FullName: "owner/repo",
			Owner:    "owner",
			Name:     "repo",
		},
	}

	terraform := mocks.NewMockClient()
	tfVersion, _ := version.NewVersion("0.15.0")
	tmpDir, cleanup := TempDir(t)
	defer cleanup()

	s := &VersionStepRunner{
		TerraformExecutor: terraform,
		DefaultTFVersion:  tfVersion,
	}

	t.Run("ensure runs", func(t *testing.T) {
		_, err := s.Run(context, []string{}, tmpDir, map[string]string(nil))
		terraform.VerifyWasCalledOnce().RunCommandWithVersion(context, tmpDir, []string{"version"}, map[string]string(nil), tfVersion, "default")
		Ok(t, err)
	})
}
