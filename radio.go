package input

import "webtyp.com/fmt"




type radio struct{ Base }

// Radio creates a new Radio input instance. opts, when given, becomes the
// closed set of choices (see Base.Validate) — the same shape Gender() uses
// for its fixed Male/Female pair, generalized to any caller-supplied list so
// a project-specific enum needs no hand-written input.Base-embedding type.
func Radio(opts ...fmt.KeyValue) Input {
	r := &radio{}
	r.Letters = true
	r.Numbers = true
	r.Minimum = 1
	r.InitBase("", "", "radio")
	r.SetOptions(opts...)
	return r
}

// SetOptions updates choices and marks the radio input FullWidth if they exceed a single-column budget.
func (r *radio) SetOptions(opts ...fmt.KeyValue) {
	r.Base.SetOptions(opts...)
	if len(opts) > 2 {
		r.FullWidth = true
	} else {
		totalChars := 0
		for _, o := range opts {
			totalChars += len(o.Value)
		}
		if totalChars > 20 {
			r.FullWidth = true
		}
	}
}

// Clone creates a new Radio input.
func (r *radio) Clone(parentID, name string) Input {
	c := *r
	c.InitBase(parentID, name, "radio")
	return &c
}
