// Code generated by pegomock. DO NOT EDIT.
package matchers

import (
	"reflect"

	"github.com/petergtz/pegomock"
	"github.com/gojekfarm/atlantis/server/events/command"
)

func AnyModelsProjectResult() command.ProjectResult {
	pegomock.RegisterMatcher(pegomock.NewAnyMatcher(reflect.TypeOf((*(command.ProjectResult))(nil)).Elem()))
	var nullValue command.ProjectResult
	return nullValue
}

func EqModelsProjectResult(value command.ProjectResult) command.ProjectResult {
	pegomock.RegisterMatcher(&pegomock.EqMatcher{Value: value})
	var nullValue command.ProjectResult
	return nullValue
}

func NotEqModelsProjectResult(value command.ProjectResult) command.ProjectResult {
	pegomock.RegisterMatcher(&pegomock.NotEqMatcher{Value: value})
	var nullValue command.ProjectResult
	return nullValue
}

func ModelsProjectResultThat(matcher pegomock.ArgumentMatcher) command.ProjectResult {
	pegomock.RegisterMatcher(matcher)
	var nullValue command.ProjectResult
	return nullValue
}
