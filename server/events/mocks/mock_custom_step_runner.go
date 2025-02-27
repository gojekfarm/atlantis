// Code generated by pegomock. DO NOT EDIT.
// Source: github.com/runatlantis/atlantis/server/events (interfaces: CustomStepRunner)

package mocks

import (
	pegomock "github.com/petergtz/pegomock/v4"
	command "github.com/runatlantis/atlantis/server/events/command"
	"reflect"
	"time"
)

type MockCustomStepRunner struct {
	fail func(message string, callerSkip ...int)
}

func NewMockCustomStepRunner(options ...pegomock.Option) *MockCustomStepRunner {
	mock := &MockCustomStepRunner{}
	for _, option := range options {
		option.Apply(mock)
	}
	return mock
}

func (mock *MockCustomStepRunner) SetFailHandler(fh pegomock.FailHandler) { mock.fail = fh }
func (mock *MockCustomStepRunner) FailHandler() pegomock.FailHandler      { return mock.fail }

func (mock *MockCustomStepRunner) Run(ctx command.ProjectContext, cmd string, path string, envs map[string]string, streamOutput bool) (string, error) {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockCustomStepRunner().")
	}
	params := []pegomock.Param{ctx, cmd, path, envs, streamOutput}
	result := pegomock.GetGenericMockFrom(mock).Invoke("Run", params, []reflect.Type{reflect.TypeOf((*string)(nil)).Elem(), reflect.TypeOf((*error)(nil)).Elem()})
	var ret0 string
	var ret1 error
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(string)
		}
		if result[1] != nil {
			ret1 = result[1].(error)
		}
	}
	return ret0, ret1
}

func (mock *MockCustomStepRunner) VerifyWasCalledOnce() *VerifierMockCustomStepRunner {
	return &VerifierMockCustomStepRunner{
		mock:                   mock,
		invocationCountMatcher: pegomock.Times(1),
	}
}

func (mock *MockCustomStepRunner) VerifyWasCalled(invocationCountMatcher pegomock.InvocationCountMatcher) *VerifierMockCustomStepRunner {
	return &VerifierMockCustomStepRunner{
		mock:                   mock,
		invocationCountMatcher: invocationCountMatcher,
	}
}

func (mock *MockCustomStepRunner) VerifyWasCalledInOrder(invocationCountMatcher pegomock.InvocationCountMatcher, inOrderContext *pegomock.InOrderContext) *VerifierMockCustomStepRunner {
	return &VerifierMockCustomStepRunner{
		mock:                   mock,
		invocationCountMatcher: invocationCountMatcher,
		inOrderContext:         inOrderContext,
	}
}

func (mock *MockCustomStepRunner) VerifyWasCalledEventually(invocationCountMatcher pegomock.InvocationCountMatcher, timeout time.Duration) *VerifierMockCustomStepRunner {
	return &VerifierMockCustomStepRunner{
		mock:                   mock,
		invocationCountMatcher: invocationCountMatcher,
		timeout:                timeout,
	}
}

type VerifierMockCustomStepRunner struct {
	mock                   *MockCustomStepRunner
	invocationCountMatcher pegomock.InvocationCountMatcher
	inOrderContext         *pegomock.InOrderContext
	timeout                time.Duration
}

func (verifier *VerifierMockCustomStepRunner) Run(ctx command.ProjectContext, cmd string, path string, envs map[string]string, streamOutput bool) *MockCustomStepRunner_Run_OngoingVerification {
	params := []pegomock.Param{ctx, cmd, path, envs, streamOutput}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "Run", params, verifier.timeout)
	return &MockCustomStepRunner_Run_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type MockCustomStepRunner_Run_OngoingVerification struct {
	mock              *MockCustomStepRunner
	methodInvocations []pegomock.MethodInvocation
}

func (c *MockCustomStepRunner_Run_OngoingVerification) GetCapturedArguments() (command.ProjectContext, string, string, map[string]string, bool) {
	ctx, cmd, path, envs, streamOutput := c.GetAllCapturedArguments()
	return ctx[len(ctx)-1], cmd[len(cmd)-1], path[len(path)-1], envs[len(envs)-1], streamOutput[len(streamOutput)-1]
}

func (c *MockCustomStepRunner_Run_OngoingVerification) GetAllCapturedArguments() (_param0 []command.ProjectContext, _param1 []string, _param2 []string, _param3 []map[string]string, _param4 []bool) {
	params := pegomock.GetGenericMockFrom(c.mock).GetInvocationParams(c.methodInvocations)
	if len(params) > 0 {
		_param0 = make([]command.ProjectContext, len(c.methodInvocations))
		for u, param := range params[0] {
			_param0[u] = param.(command.ProjectContext)
		}
		_param1 = make([]string, len(c.methodInvocations))
		for u, param := range params[1] {
			_param1[u] = param.(string)
		}
		_param2 = make([]string, len(c.methodInvocations))
		for u, param := range params[2] {
			_param2[u] = param.(string)
		}
		_param3 = make([]map[string]string, len(c.methodInvocations))
		for u, param := range params[3] {
			_param3[u] = param.(map[string]string)
		}
		_param4 = make([]bool, len(c.methodInvocations))
		for u, param := range params[4] {
			_param4[u] = param.(bool)
		}
	}
	return
}
