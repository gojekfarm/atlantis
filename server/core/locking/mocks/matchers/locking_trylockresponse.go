// Code generated by pegomock. DO NOT EDIT.
package matchers

import (
	"github.com/petergtz/pegomock"
	"reflect"

	locking "github.com/gojekfarm/atlantis/server/core/locking"
)

func AnyLockingTryLockResponse() locking.TryLockResponse {
	pegomock.RegisterMatcher(pegomock.NewAnyMatcher(reflect.TypeOf((*(locking.TryLockResponse))(nil)).Elem()))
	var nullValue locking.TryLockResponse
	return nullValue
}

func EqLockingTryLockResponse(value locking.TryLockResponse) locking.TryLockResponse {
	pegomock.RegisterMatcher(&pegomock.EqMatcher{Value: value})
	var nullValue locking.TryLockResponse
	return nullValue
}

func NotEqLockingTryLockResponse(value locking.TryLockResponse) locking.TryLockResponse {
	pegomock.RegisterMatcher(&pegomock.NotEqMatcher{Value: value})
	var nullValue locking.TryLockResponse
	return nullValue
}

func LockingTryLockResponseThat(matcher pegomock.ArgumentMatcher) locking.TryLockResponse {
	pegomock.RegisterMatcher(matcher)
	var nullValue locking.TryLockResponse
	return nullValue
}
