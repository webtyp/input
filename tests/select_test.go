package input_test

import (
	"testing"

	"webtyp.com/fmt"
	"webtyp.com/input"
)

// A select whose options arrive at runtime holds record ids, which carry '-'
// and '_'; HTML-dangerous characters stay rejected.
func TestSelectWithoutOptionsAcceptsIDs(t *testing.T) {
	for _, v := range []string{"test-id-1", "zone_2", "1789255672825978415"} {
		if err := input.Select().Validate(v); err != nil {
			t.Errorf("Select().Validate(%q) = %v, want nil", v, err)
		}
	}
	for _, v := range []string{"a<b", "x\"y", "a b"} {
		if err := input.Select().Validate(v); err == nil {
			t.Errorf("Select().Validate(%q) = nil, want an error", v)
		}
	}
}

// With options, only membership counts.
func TestSelectWithOptionsChecksMembership(t *testing.T) {
	s := input.Select(fmt.KeyValue{Key: "internet_filtered", Value: "Internet (filtered)"})
	if err := s.Validate("internet_filtered"); err != nil {
		t.Errorf("member rejected: %v", err)
	}
	if err := s.Validate("test-id-1"); err == nil {
		t.Error("non-member accepted")
	}
}
