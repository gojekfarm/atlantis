// Code generated by pegomock. DO NOT EDIT.
// Source: github.com/runatlantis/atlantis/server/events (interfaces: JobURLSetter)

package mocks

import (
	"reflect"
	"time"

	pegomock "github.com/petergtz/pegomock"
	"github.com/gojekfarm/atlantis/server/events/command"
	models "github.com/gojekfarm/atlantis/server/events/models"
)

type MockJobURLSetter struct {
	fail func(message string, callerSkip ...int)
}

func NewMockJobURLSetter(options ...pegomock.Option) *MockJobURLSetter {
	mock := &MockJobURLSetter{}
	for _, option := range options {
		option.Apply(mock)
	}
	return mock
}

func (mock *MockJobURLSetter) SetFailHandler(fh pegomock.FailHandler) { mock.fail = fh }
func (mock *MockJobURLSetter) FailHandler() pegomock.FailHandler      { return mock.fail }

func (mock *MockJobURLSetter) SetJobURLWithStatus(_param0 command.ProjectContext, _param1 command.Name, _param2 models.CommitStatus) error {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockJobURLSetter().")
	}
	params := []pegomock.Param{_param0, _param1, _param2}
	result := pegomock.GetGenericMockFrom(mock).Invoke("SetJobURLWithStatus", params, []reflect.Type{reflect.TypeOf((*error)(nil)).Elem()})
	var ret0 error
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(error)
		}
	}
	return ret0
}

func (mock *MockJobURLSetter) VerifyWasCalledOnce() *VerifierMockJobURLSetter {
	return &VerifierMockJobURLSetter{
		mock:                   mock,
		invocationCountMatcher: pegomock.Times(1),
	}
}

func (mock *MockJobURLSetter) VerifyWasCalled(invocationCountMatcher pegomock.InvocationCountMatcher) *VerifierMockJobURLSetter {
	return &VerifierMockJobURLSetter{
		mock:                   mock,
		invocationCountMatcher: invocationCountMatcher,
	}
}

func (mock *MockJobURLSetter) VerifyWasCalledInOrder(invocationCountMatcher pegomock.InvocationCountMatcher, inOrderContext *pegomock.InOrderContext) *VerifierMockJobURLSetter {
	return &VerifierMockJobURLSetter{
		mock:                   mock,
		invocationCountMatcher: invocationCountMatcher,
		inOrderContext:         inOrderContext,
	}
}

func (mock *MockJobURLSetter) VerifyWasCalledEventually(invocationCountMatcher pegomock.InvocationCountMatcher, timeout time.Duration) *VerifierMockJobURLSetter {
	return &VerifierMockJobURLSetter{
		mock:                   mock,
		invocationCountMatcher: invocationCountMatcher,
		timeout:                timeout,
	}
}

type VerifierMockJobURLSetter struct {
	mock                   *MockJobURLSetter
	invocationCountMatcher pegomock.InvocationCountMatcher
	inOrderContext         *pegomock.InOrderContext
	timeout                time.Duration
}

func (verifier *VerifierMockJobURLSetter) SetJobURLWithStatus(_param0 command.ProjectContext, _param1 command.Name, _param2 models.CommitStatus) *MockJobURLSetter_SetJobURLWithStatus_OngoingVerification {
	params := []pegomock.Param{_param0, _param1, _param2}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "SetJobURLWithStatus", params, verifier.timeout)
	return &MockJobURLSetter_SetJobURLWithStatus_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type MockJobURLSetter_SetJobURLWithStatus_OngoingVerification struct {
	mock              *MockJobURLSetter
	methodInvocations []pegomock.MethodInvocation
}

func (c *MockJobURLSetter_SetJobURLWithStatus_OngoingVerification) GetCapturedArguments() (command.ProjectContext, command.Name, models.CommitStatus) {
	_param0, _param1, _param2 := c.GetAllCapturedArguments()
	return _param0[len(_param0)-1], _param1[len(_param1)-1], _param2[len(_param2)-1]
}

func (c *MockJobURLSetter_SetJobURLWithStatus_OngoingVerification) GetAllCapturedArguments() (_param0 []command.ProjectContext, _param1 []command.Name, _param2 []models.CommitStatus) {
	params := pegomock.GetGenericMockFrom(c.mock).GetInvocationParams(c.methodInvocations)
	if len(params) > 0 {
		_param0 = make([]command.ProjectContext, len(c.methodInvocations))
		for u, param := range params[0] {
			_param0[u] = param.(command.ProjectContext)
		}
		_param1 = make([]command.Name, len(c.methodInvocations))
		for u, param := range params[1] {
			_param1[u] = param.(command.Name)
		}
		_param2 = make([]models.CommitStatus, len(c.methodInvocations))
		for u, param := range params[2] {
			_param2[u] = param.(models.CommitStatus)
		}
	}
	return
}
