// Code generated by pegomock. DO NOT EDIT.
package matchers

import (
	"reflect"

	"github.com/petergtz/pegomock"

	events "github.com/gojekfarm/atlantis/server/events"
)

func AnyPtrToEventsTryLockResponse() *events.TryLockResponse {
	pegomock.RegisterMatcher(pegomock.NewAnyMatcher(reflect.TypeOf((*(*events.TryLockResponse))(nil)).Elem()))
	var nullValue *events.TryLockResponse
	return nullValue
}

func EqPtrToEventsTryLockResponse(value *events.TryLockResponse) *events.TryLockResponse {
	pegomock.RegisterMatcher(&pegomock.EqMatcher{Value: value})
	var nullValue *events.TryLockResponse
	return nullValue
}

func NotEqPtrToEventsTryLockResponse(value *events.TryLockResponse) *events.TryLockResponse {
	pegomock.RegisterMatcher(&pegomock.NotEqMatcher{Value: value})
	var nullValue *events.TryLockResponse
	return nullValue
}

func PtrToEventsTryLockResponseThat(matcher pegomock.ArgumentMatcher) *events.TryLockResponse {
	pegomock.RegisterMatcher(matcher)
	var nullValue *events.TryLockResponse
	return nullValue
}
