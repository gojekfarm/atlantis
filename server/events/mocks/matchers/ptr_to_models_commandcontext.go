// Code generated by pegomock. DO NOT EDIT.
package matchers

import (
	"reflect"

	"github.com/petergtz/pegomock"
	"github.com/gojekfarm/atlantis/server/events/command"
)

func AnyPtrToModelsCommandContext() *command.Context {
	pegomock.RegisterMatcher(pegomock.NewAnyMatcher(reflect.TypeOf((*(*command.Context))(nil)).Elem()))
	var nullValue *command.Context
	return nullValue
}

func EqPtrToModelsCommandContext(value *command.Context) *command.Context {
	pegomock.RegisterMatcher(&pegomock.EqMatcher{Value: value})
	var nullValue *command.Context
	return nullValue
}
