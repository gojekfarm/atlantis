// Code generated by pegomock. DO NOT EDIT.
package matchers

import (
	"github.com/petergtz/pegomock"
	"reflect"

	models "github.com/gojekfarm/atlantis/server/events/models"
)

func AnyModelsProject() models.Project {
	pegomock.RegisterMatcher(pegomock.NewAnyMatcher(reflect.TypeOf((*(models.Project))(nil)).Elem()))
	var nullValue models.Project
	return nullValue
}

func EqModelsProject(value models.Project) models.Project {
	pegomock.RegisterMatcher(&pegomock.EqMatcher{Value: value})
	var nullValue models.Project
	return nullValue
}

func NotEqModelsProject(value models.Project) models.Project {
	pegomock.RegisterMatcher(&pegomock.NotEqMatcher{Value: value})
	var nullValue models.Project
	return nullValue
}

func ModelsProjectThat(matcher pegomock.ArgumentMatcher) models.Project {
	pegomock.RegisterMatcher(matcher)
	var nullValue models.Project
	return nullValue
}
