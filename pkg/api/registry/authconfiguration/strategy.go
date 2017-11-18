package auth

import (
	"fmt"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"

	"kope.io/auth/pkg/apis/componentconfig"
)

type apiServerStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

func NewStrategy(typer runtime.ObjectTyper) apiServerStrategy {
	return apiServerStrategy{typer, names.SimpleNameGenerator}
}

func (apiServerStrategy) NamespaceScoped() bool {
	return false
}

func (apiServerStrategy) PrepareForCreate(ctx genericapirequest.Context, obj runtime.Object) {
}

func (apiServerStrategy) PrepareForUpdate(ctx genericapirequest.Context, obj, old runtime.Object) {
}

func (apiServerStrategy) Validate(ctx genericapirequest.Context, obj runtime.Object) field.ErrorList {
	return field.ErrorList{}
	// return validation.ValidateFlunder(obj.(*wardle.Flunder))
}

func (apiServerStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (apiServerStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (apiServerStrategy) Canonicalize(obj runtime.Object) {
}

func (apiServerStrategy) ValidateUpdate(ctx genericapirequest.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
	// return validation.ValidateFlunderUpdate(obj.(*wardle.Flunder), old.(*wardle.Flunder))
}

func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, bool, error) {
	o, ok := obj.(*componentconfig.AuthConfiguration)
	if !ok {
		return nil, nil, false, fmt.Errorf("given object is not a AuthConfiguration.")
	}
	return labels.Set(o.ObjectMeta.Labels), AuthConfigurationToSelectableFields(o), o.Initializers != nil, nil
}

func MatchAuthConfiguration(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

func AuthConfigurationToSelectableFields(obj *componentconfig.AuthConfiguration) fields.Set {
	return generic.ObjectMetaFieldsSet(&obj.ObjectMeta, true)
}
