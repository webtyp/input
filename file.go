package input

import (
	"webtyp.com/filetype"
	"webtyp.com/fmt"
)

// FileKind is the kind of a field that holds a file. The record keeps only the file's key (a text column); the
// bytes are validated by ValidateFile, in the browser before sending and in crud on arrival: the same code.
type FileKind interface {
	Input
	// ValidateFile reads the real type from the bytes (never from what the client declared) and refuses what is
	// not in the allowed list or is larger than MaxBytes.
	ValidateFile(data []byte) (filetype.Type, error)
	// MaxBytes is the largest file accepted.
	MaxBytes() int
}

type file struct {
	Base
	allowed  filetype.Allowlist
	maxBytes int
}

// File creates the kind of a file field. allowed lists the accepted types (filetype.Images for raster images; never
// SVG or HTML: they run JavaScript); maxBytes is the largest size.
//
//	logo := input.File(filetype.Images, 2<<20)
func File(allowed filetype.Allowlist, maxBytes int) FileKind {
	f := &file{allowed: allowed, maxBytes: maxBytes}
	f.Letters = true
	f.Numbers = true
	f.Extra = []rune{'-', '_', '.'} // the key the server mints: an id plus an extension; never a path
	f.Minimum = 1
	f.Maximum = 128
	f.InitBase("", "", "file")
	return f
}

// Clone creates a new File input with the given parentID and name, keeping the allowed types and the size.
func (f *file) Clone(parentID, name string) Input {
	c := *f
	c.InitBase(parentID, name, "file")
	return &c
}

func (f *file) MaxBytes() int { return f.maxBytes }

func (f *file) ValidateFile(data []byte) (filetype.Type, error) {
	if f.maxBytes > 0 && len(data) > f.maxBytes {
		return filetype.Type{}, fmt.Err("File", "Large")
	}
	return f.allowed.Validate(data)
}
