package input

import "webtyp.com/fmt"

type rut struct{ Base }

// Rut creates a new RUT input instance.
func Rut() Input {
	r := &rut{}
	r.Numbers = true
	r.Extra = []rune{'-', 'k', 'K', '.', ' '}
	r.Minimum = 3
	r.Maximum = 16 // " 12.345.678 - 9 "
	r.InitBase("", "", "text")
	r.SetPlaceholder("example:")
	r.SetExample("12345678-9")
	return r
}

// Validate validates the Chilean RUT format and check digit (see NormalizeRUT: any way of writing it is accepted).
func (r *rut) Validate(value string) error {
	if err := r.Permitted.Validate(r.name, value); err != nil {
		return err
	}
	_, err := NormalizeRUT(value)
	return err
}

// NormalizeRUT is the only RUT validator of the ecosystem: it checks the check digit and returns the RUT in its
// one canonical form, "body-DV" with an upper-case K. Every way a person writes the same RUT gives the same
// result — dots, hyphen, any space (also inside) and leading zeros are dropped — so a lookup by RUT never misses
// a record that already exists.
//
// Example: norm, err := input.NormalizeRUT(" 08.123.456 - 6 ") // "8123456-6", nil
func NormalizeRUT(value string) (string, error) {
	clean := make([]byte, 0, len(value))
	for i := 0; i < len(value); i++ {
		switch c := value[i]; c {
		case '.', '-', ' ', '\t', '\n', '\r':
		default:
			if c == '0' && len(clean) == 0 {
				continue // leading zero: "08123456-6" is "8123456-6"
			}
			clean = append(clean, c)
		}
	}
	if len(clean) < 2 {
		return "", fmt.Err("Format", "Invalid")
	}
	body, dv := string(clean[:len(clean)-1]), string(clean[len(clean)-1:])
	if dv == "k" {
		dv = "K"
	}
	for i := 0; i < len(body); i++ {
		if body[i] < '0' || body[i] > '9' {
			return "", fmt.Err("Format", "Invalid")
		}
	}
	if RutCheckDigit(body) != dv {
		return "", fmt.Err("Digit", "Invalid")
	}
	return body + "-" + dv, nil
}

// RutCheckDigit returns the check digit ("0"-"9" or "K") of a RUT body made only of digits (modulo 11).
//
// Example: input.RutCheckDigit("8123456") // "6"
func RutCheckDigit(body string) string {
	sum, factor := 0, 2
	for i := len(body) - 1; i >= 0; i-- {
		sum += int(body[i]-'0') * factor
		if factor == 7 {
			factor = 2
		} else {
			factor++
		}
	}
	switch val := 11 - sum%11; val {
	case 11:
		return "0"
	case 10:
		return "K"
	default:
		return fmt.Convert(val).String()
	}
}

// Clone creates a new rut input with the given parentID and name.
func (r *rut) Clone(parentID, name string) Input {
	c := *r
	c.InitBase(parentID, name, "text")
	return &c
}
