// Code generated by pegomock. DO NOT EDIT.
package matchers

import (
	"github.com/petergtz/pegomock"
	"reflect"

	command "github.com/gojekfarm/atlantis/server/events/command"
)

func AnyCommandName() command.Name {
	pegomock.RegisterMatcher(pegomock.NewAnyMatcher(reflect.TypeOf((*(command.Name))(nil)).Elem()))
	var nullValue command.Name
	return nullValue
}

func EqCommandName(value command.Name) command.Name {
	pegomock.RegisterMatcher(&pegomock.EqMatcher{Value: value})
	var nullValue command.Name
	return nullValue
}

func NotEqCommandName(value command.Name) command.Name {
	pegomock.RegisterMatcher(&pegomock.NotEqMatcher{Value: value})
	var nullValue command.Name
	return nullValue
}

func CommandNameThat(matcher pegomock.ArgumentMatcher) command.Name {
	pegomock.RegisterMatcher(matcher)
	var nullValue command.Name
	return nullValue
}
