package controllers_test

import (
	"bytes"
	"encoding/json"
	"github.com/gojekfarm/atlantis/server/events/command"
	"github.com/gojekfarm/atlantis/server/events/models"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/petergtz/pegomock"
	"github.com/gojekfarm/atlantis/server/controllers"
	. "github.com/gojekfarm/atlantis/server/core/locking/mocks"
	"github.com/gojekfarm/atlantis/server/events"
	. "github.com/gojekfarm/atlantis/server/events/mocks"
	. "github.com/gojekfarm/atlantis/server/events/mocks/matchers"
	. "github.com/gojekfarm/atlantis/server/events/vcs/mocks"
	"github.com/gojekfarm/atlantis/server/logging"
	"github.com/gojekfarm/atlantis/server/metrics"
	. "github.com/gojekfarm/atlantis/testing"
)

const atlantisTokenHeader = "X-Atlantis-Token"
const atlantisToken = "token"

func TestAPIController_Plan(t *testing.T) {
	ac, projectCommandBuilder, projectCommandRunner := setup(t)
	body, _ := json.Marshal(controllers.APIRequest{
		Repository: "Repo",
		Ref:        "main",
		Type:       "Gitlab",
		Projects:   []string{"default"},
	})
	req, _ := http.NewRequest("POST", "", bytes.NewBuffer(body))
	req.Header.Set(atlantisTokenHeader, atlantisToken)
	w := httptest.NewRecorder()
	ac.Plan(w, req)
	ResponseContains(t, w, http.StatusOK, "")
	projectCommandBuilder.VerifyWasCalledOnce().BuildPlanCommands(AnyPtrToEventsCommandContext(), AnyPtrToEventsCommentCommand())
	projectCommandRunner.VerifyWasCalledOnce().Plan(AnyModelsProjectCommandContext())
}

func TestAPIController_Apply(t *testing.T) {
	ac, projectCommandBuilder, projectCommandRunner := setup(t)
	body, _ := json.Marshal(controllers.APIRequest{
		Repository: "Repo",
		Ref:        "main",
		Type:       "Gitlab",
		Projects:   []string{"default"},
	})
	req, _ := http.NewRequest("POST", "", bytes.NewBuffer(body))
	req.Header.Set(atlantisTokenHeader, atlantisToken)
	w := httptest.NewRecorder()
	ac.Apply(w, req)
	ResponseContains(t, w, http.StatusOK, "")
	projectCommandBuilder.VerifyWasCalledOnce().BuildApplyCommands(AnyPtrToEventsCommandContext(), AnyPtrToEventsCommentCommand())
	projectCommandRunner.VerifyWasCalledOnce().Plan(AnyModelsProjectCommandContext())
	projectCommandRunner.VerifyWasCalledOnce().Apply(AnyModelsProjectCommandContext())
}

func setup(t *testing.T) (controllers.APIController, *MockProjectCommandBuilder, *MockProjectCommandRunner) {
	RegisterMockTestingT(t)
	locker := NewMockLocker()
	logger := logging.NewNoopLogger(t)
	scope, _, _ := metrics.NewLoggingScope(logger, "null")
	parser := NewMockEventParsing()
	vcsClient := NewMockClient()
	repoAllowlistChecker, err := events.NewRepoAllowlistChecker("*")
	Ok(t, err)

	projectCommandBuilder := NewMockProjectCommandBuilder()
	When(projectCommandBuilder.BuildPlanCommands(AnyPtrToEventsCommandContext(), AnyPtrToEventsCommentCommand())).
		ThenReturn([]command.ProjectContext{{
			CommandName: command.Plan,
		}}, nil)
	When(projectCommandBuilder.BuildApplyCommands(AnyPtrToEventsCommandContext(), AnyPtrToEventsCommentCommand())).
		ThenReturn([]command.ProjectContext{{
			CommandName: command.Apply,
		}}, nil)

	projectCommandRunner := NewMockProjectCommandRunner()
	When(projectCommandRunner.Plan(AnyModelsProjectCommandContext())).ThenReturn(command.ProjectResult{
		PlanSuccess: &models.PlanSuccess{},
	})
	When(projectCommandRunner.Apply(AnyModelsProjectCommandContext())).ThenReturn(command.ProjectResult{
		ApplySuccess: "success",
	})

	ac := controllers.APIController{
		APISecret:                 []byte(atlantisToken),
		Locker:                    locker,
		Logger:                    logger,
		Scope:                     scope,
		Parser:                    parser,
		ProjectCommandBuilder:     projectCommandBuilder,
		ProjectPlanCommandRunner:  projectCommandRunner,
		ProjectApplyCommandRunner: projectCommandRunner,
		VCSClient:                 vcsClient,
		RepoAllowlistChecker:      repoAllowlistChecker,
	}
	return ac, projectCommandBuilder, projectCommandRunner
}
