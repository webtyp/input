package input_test

import (
	"strings"
	"testing"
	"webtyp.com/filetype"
	"webtyp.com/input"
	"webtyp.com/model"
)

var (
	png = []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A, 0, 0}
	svgBytes = []byte("<svg xmlns='http://www.w3.org/2000/svg'></svg>")
)

func TestFile_AcceptsAllowedTypeFromBytes(t *testing.T) {
	k := input.File(filetype.Images, 1024)
	typ, err := k.ValidateFile(png)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if typ != filetype.PNG {
		t.Errorf("expected type PNG, got %v", typ)
	}
}

func TestFile_IgnoresWhatTheClientClaims(t *testing.T) {
	k := input.File(filetype.NewAllowlist(filetype.PNG), 1024)
	_, err := k.ValidateFile(svgBytes)
	if err == nil {
		t.Errorf("expected error when bytes do not match allowlist, got nil")
	}
}

func TestFile_RejectsSVG(t *testing.T) {
	k := input.File(filetype.Images, 1024)
	_, err := k.ValidateFile(svgBytes)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "SVG") && !strings.Contains(err.Error(), "svg") {
		t.Errorf("expected error whose text contains SVG, got %v", err)
	}
}

func TestFile_RejectsOverMaxBytes(t *testing.T) {
	k := input.File(filetype.Images, 4)
	_, err := k.ValidateFile(png)
	if err == nil {
		t.Fatalf("expected error for size over max bytes, got nil")
	}
	if k.MaxBytes() != 4 {
		t.Errorf("expected MaxBytes 4, got %v", k.MaxBytes())
	}
}

func TestFile_KeyValidation(t *testing.T) {
	k := input.File(filetype.Images, 1024)
	err := k.Validate("01HX9-a_b.png")
	if err != nil {
		t.Errorf("expected no error for valid key, got %v", err)
	}
	err = k.Validate("../x.png")
	if err == nil {
		t.Errorf("expected error for key with .., got nil")
	}
	err = k.Validate("a/b.png")
	if err == nil {
		t.Errorf("expected error for key with /, got nil")
	}

	k.SetRequired(true)
	err = k.Validate("")
	if err == nil {
		t.Errorf("expected error for empty key when required, got nil")
	}
}

func TestFile_StorageIsText(t *testing.T) {
	k := input.File(filetype.Images, 1024)
	if k.Storage() != model.FieldText {
		t.Errorf("expected storage FieldText, got %v", k.Storage())
	}
	if k.Name() != "file" {
		t.Errorf("expected html name file, got %v", k.Name())
	}
}

func TestFile_CloneKeepsTheRules(t *testing.T) {
	k := input.File(filetype.Images, 1024)
	c := k.Clone("form", "logo").(input.FileKind)
	if c.MaxBytes() != k.MaxBytes() {
		t.Errorf("expected cloned MaxBytes to be %v, got %v", k.MaxBytes(), c.MaxBytes())
	}
	_, err := c.ValidateFile(svgBytes)
	if err == nil {
		t.Errorf("expected error when validating svg on cloned input, got nil")
	}
}
