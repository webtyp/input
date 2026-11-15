package input
import "webtyp.com/model"


import (
	"webtyp.com/fmt"
)

// Input interface defines the behavior for all form input types.
type Input interface {
	model.Kind
	Clone(parentID, name string) Input
	FieldName() string
	SetRequired(bool)
	AddAttribute(key, value string)

	// Metadata getters for rendering
	GetID() string
	SetID(string)
	GetValues() []string
	GetOptions() []fmt.KeyValue
	GetPlaceholder() string
	GetTitle() string
	IsRequired() bool
	IsDisabled() bool
	IsReadonly() bool
	GetAttributes() []fmt.KeyValue
	ErrorID() string
	HTMLName() string
	HandlerName() string
}

// AlsoSetter is an optional capability for an input that, when the person
// commits its value, also sets another field of the form: the one the model
// names in Field.AlsoSetsField. The form hands the input the function that
// writes that field (BindAlsoSets); the input decides WHEN to call it and WITH
// WHAT, so neither the form nor the module knows what an age or a country
// means. A value written this way follows the target input's own format
// ("true"/"false" for a checkbox).
type AlsoSetter interface {
	BindAlsoSets(set func(value string))
}
