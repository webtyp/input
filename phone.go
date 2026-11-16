package input

import "webtyp.com/country"

// phone is the one phone input: its value is always international ("+56911111111"),
// with the digits the country of its prefix has (country.Validate, the single
// rule). The person does not type the prefix: form draws a select of countries
// and the number (see form's phone control); a form with no countries shows one
// box where the number is typed with its "+".
type phone struct{ Base }

// Phone creates a new phone input instance.
func Phone() Input {
	p := &phone{}
	p.Numbers = true
	p.Extra = []rune{'+'}
	p.Minimum = 8
	p.Maximum = 16
	p.InitBase("", "", "tel")
	return p
}

// Validate accepts an international number (country.Validate); empty follows the
// base rule (an optional field is skipped by the form and by crud before it gets here).
func (p *phone) Validate(value string) error {
	if value == "" {
		return p.Base.Validate(value)
	}
	return country.Validate(value)
}

// Clone creates a new phone input with the given parentID and name.
func (p *phone) Clone(parentID, name string) Input {
	c := *p
	c.InitBase(parentID, name, "tel")
	c.SetOptions(p.GetOptions()...)
	return &c
}

// SearchableText: a phone enters the row's text_search, so «9 1234 5678» finds it (model.SearchableText).
func (p *phone) SearchableText() {}
