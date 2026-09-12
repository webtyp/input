package input

import "webtyp.com/fmt"




type select_ struct{ Base }

// Select creates a new Select input instance. opts, when given, becomes the
// closed set of choices rendered as <option> elements (see Base.Validate) —
// so a project-specific dropdown needs no hand-written input.Base-embedding
// type, the same way Radio(opts...) covers the radio-group case.
func Select(opts ...fmt.KeyValue) Input {
	s := &select_{}
	s.Letters = true
	s.Numbers = true
	s.Minimum = 1
	s.InitBase("", "", "select")
	s.SetOptions(opts...)
	return s
}

// Clone creates a new Select input.
func (s *select_) Clone(parentID, name string) Input {
	c := *s
	c.InitBase(parentID, name, "select")
	return &c
}
