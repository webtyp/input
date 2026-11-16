package input_test

import (
	"testing"

	"webtyp.com/input"
	"webtyp.com/model"
)

func TestSearchableText_ShortFreeTextKindsEnter(t *testing.T) {
	for name, k := range map[string]input.Input{"text": input.Text(), "phone": input.Phone(), "email": input.Email(), "rut": input.Rut()} {
		if _, ok := k.(model.SearchableText); !ok {
			t.Errorf("%s must implement model.SearchableText", name)
		}
	}
}

func TestSearchableText_OtherKindsStayOut(t *testing.T) {
	for name, k := range map[string]input.Input{"textarea": input.Textarea(), "address": input.Address(), "select": input.Select(), "checkbox": input.Checkbox(), "gender": input.Gender(), "password": input.Password()} {
		if _, ok := k.(model.SearchableText); ok {
			t.Errorf("%s must not implement model.SearchableText", name)
		}
	}
}
