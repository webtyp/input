package input

import "testing"

func TestMasked_DefaultsFalse(t *testing.T) {
	if Text().(interface{ IsMasked() bool }).IsMasked() {
		t.Error("Text() must not be masked by default")
	}
}

func TestPassword_IsMaskedByDefault(t *testing.T) {
	if !Password().(interface{ IsMasked() bool }).IsMasked() {
		t.Error("Password() must be masked by default")
	}
}

func TestMasked_SurvivesClone(t *testing.T) {
	proto := Text()
	proto.(interface{ SetMasked(bool) }).SetMasked(true)

	cloned := proto.Clone("parent", "field")
	if !cloned.(interface{ IsMasked() bool }).IsMasked() {
		t.Error("SetMasked(true) must survive Clone")
	}
}

func TestMasked_DoesNotAffectValidation(t *testing.T) {
	inp := Text()
	inp.(interface{ SetMasked(bool) }).SetMasked(true)

	// A masked Text() still validates exactly like an unmasked one —
	// masking is presentation, not a character-set or checksum change.
	if err := inp.Validate("ok value"); err != nil {
		t.Errorf("masked Text() rejected a value its own charset allows: %v", err)
	}
	if err := inp.Validate("bad:value"); err == nil {
		t.Error("masked Text() must still reject a character its charset disallows")
	}
}
