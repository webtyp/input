package input_test

import (
	"testing"

	"webtyp.com/input"
)

const canonical = "48:F1:7F:D9:D7:B7"

var validMACs = []string{
	"48:F1:7F:D9:D7:B7",
	"48-f1-7f-d9-d7-b7",
	"48F17FD9D7B7",
}

// A form value with surrounding space is rejected by Validate (Permitted has
// no Spaces, same as IP()), but CanonicalMAC trims it: lookups of a value read
// from elsewhere (a router, a pasted ipconfig line) still match.
const spacedMAC = " 48:f1:7f:d9:d7:b7 "

func TestMACValidateAccepts(t *testing.T) {
	for _, v := range validMACs {
		if err := input.MAC().Validate(v); err != nil {
			t.Errorf("Validate(%q) = %v, want nil", v, err)
		}
	}
}

func TestMACValidateRejects(t *testing.T) {
	invalid := []string{
		"48:F1:7F:D9:D7",       // 5 groups
		"48:F1:7F:D9:D7:B7:00", // 7 groups
		"48:F1-7F:D9:D7:B7",    // mixed separators
		"GG:F1:7F:D9:D7:B7",    // not hex
		"48F17FD9D7B",          // 11 digits
		"48.F1.7F.D9.D7.B7",    // unsupported separator
		spacedMAC,              // surrounding space (same rule as IP)
	}
	for _, v := range invalid {
		if err := input.MAC().Validate(v); err == nil {
			t.Errorf("Validate(%q) = nil, want an error", v)
		}
	}
}

// An empty value behaves exactly like IP().Validate(""): both kinds put the
// shortest valid spelling in Minimum, so Permitted rejects or accepts "" the
// same way for both.
func TestMACValidateEmptyMatchesIP(t *testing.T) {
	macErr := input.MAC().Validate("")
	ipErr := input.IP().Validate("")
	if (macErr == nil) != (ipErr == nil) {
		t.Errorf("MAC().Validate(\"\") = %v, IP().Validate(\"\") = %v: want the same outcome", macErr, ipErr)
	}
}

func TestCanonicalMAC(t *testing.T) {
	for _, v := range validMACs {
		if got := input.CanonicalMAC(v); got != canonical {
			t.Errorf("CanonicalMAC(%q) = %q, want %q", v, got, canonical)
		}
	}
	if got := input.CanonicalMAC(spacedMAC); got != canonical {
		t.Errorf("CanonicalMAC(%q) = %q, want %q", spacedMAC, got, canonical)
	}
	if got := input.CanonicalMAC(" not-a-mac "); got != "NOT-A-MAC" {
		t.Errorf("CanonicalMAC of a non-MAC = %q, want trimmed upper-case %q", got, "NOT-A-MAC")
	}
}

func TestIsLocallyAdministeredMAC(t *testing.T) {
	// Real addresses from the incident of 2026-10-07: Windows "change daily",
	// an iPhone private address, a tablet, a phone.
	randomized := []string{"4A:C5:93:7A:12:DE", "DE:94:42:52:C5:59", "12:AB:2F:3F:6C:A3", "CE:45:90:5F:05:E7", "4a-c5-93-7a-12-de"}
	for _, v := range randomized {
		if !input.IsLocallyAdministeredMAC(v) {
			t.Errorf("IsLocallyAdministeredMAC(%q) = false, want true", v)
		}
	}
	real := []string{"48:F1:7F:D9:D7:B7", "60:6C:66:C0:3F:1E", "not-a-mac"}
	for _, v := range real {
		if input.IsLocallyAdministeredMAC(v) {
			t.Errorf("IsLocallyAdministeredMAC(%q) = true, want false", v)
		}
	}
}

// Consumer-shaped: the field a model.Definition declares is cloned per form.
func TestMACClonedField(t *testing.T) {
	field := input.MAC().Clone("form1", "mac")
	if err := field.Validate("48-F1-7F-D9-D7-B7"); err != nil {
		t.Errorf("cloned field rejected a valid MAC: %v", err)
	}
	if err := field.Validate("48:F1"); err == nil {
		t.Error("cloned field accepted \"48:F1\"")
	}
}
