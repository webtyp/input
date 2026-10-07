package input

import "webtyp.com/fmt"

type mac struct{ Base }

// macSeparator is the separator CanonicalMAC writes.
const macSeparator = ':'

// macHexDigits is the number of hex digits in a MAC address (6 bytes).
const macHexDigits = 12

// macGroups is the number of 2-digit groups in a separated MAC address.
const macGroups = 6

// macLocallyAdministeredBit is bit 0x02 of the first byte (IEEE 802).
const macLocallyAdministeredBit = 0x2

// MAC creates a MAC address input. It accepts "AA:BB:CC:DD:EE:FF",
// "AA-BB-CC-DD-EE-FF" and "AABBCCDDEEFF", in any case.
func MAC() Input {
	i := &mac{}
	i.Numbers = true
	i.Letters = true // hex digits a-f
	i.Extra = []rune{':', '-'}
	i.Minimum = macHexDigits                 // "AABBCCDDEEFF"
	i.Maximum = macHexDigits + macGroups - 1 // "AA:BB:CC:DD:EE:FF"
	i.InitBase("", "", "text")
	// The three accepted spellings cannot be inferred from the field name —
	// the format hint belongs to the type that defines the format (ip.go).
	i.SetPlaceholder("example:", "48:F1:7F:D9:D7:B7")
	return i
}

// Validate validates the MAC address format.
func (i *mac) Validate(value string) error {
	if err := i.Permitted.Validate(i.name, value); err != nil {
		return err
	}
	if _, ok := parseMAC(value); !ok {
		return fmt.Err("Format", "Invalid")
	}
	return nil
}

// Clone satisfies input.Input — MAC() returns Input which implements it.
func (i *mac) Clone(parentID, name string) Input {
	c := *i
	c.InitBase(parentID, name, "text")
	return &c
}

// CanonicalMAC returns the single spelling two MACs are compared by:
// upper case, colon-separated ("48:F1:7F:D9:D7:B7"). Whoever stores a MAC
// and whoever later looks it up must both pass it through here.
// A value that is not a MAC is returned trimmed and upper-cased, unchanged
// otherwise — MAC().Validate is the validator.
func CanonicalMAC(value string) string {
	digits, ok := parseMAC(value)
	if !ok {
		return fmt.Convert(value).TrimSpace().ToUpper().String()
	}
	out := make([]byte, 0, macHexDigits+macGroups-1)
	for g := 0; g < macGroups; g++ {
		if g > 0 {
			out = append(out, macSeparator)
		}
		out = append(out, digits[2*g], digits[2*g+1])
	}
	return string(out)
}

// IsLocallyAdministeredMAC reports whether the MAC has the IEEE "locally
// administered" bit set (second hex digit 2, 6, A or E). That is the case
// for randomized/private addresses that phones and Windows invent per
// Wi-Fi network, and for virtual machines. False for an invalid MAC.
func IsLocallyAdministeredMAC(value string) bool {
	digits, ok := parseMAC(value)
	if !ok {
		return false
	}
	return hexValue(digits[1])&macLocallyAdministeredBit != 0
}

// parseMAC returns the 12 hex digits of value in upper case, or ok=false.
// Accepted shapes: 12 hex digits; or 6 groups of 2 hex digits separated by
// ':' everywhere or by '-' everywhere (mixed separators are rejected).
func parseMAC(value string) (digits string, ok bool) {
	v := fmt.Convert(value).TrimSpace().String()
	out := make([]byte, 0, macHexDigits)

	switch len(v) {
	case macHexDigits:
		for k := 0; k < len(v); k++ {
			c, valid := upperHex(v[k])
			if !valid {
				return "", false
			}
			out = append(out, c)
		}
	case macHexDigits + macGroups - 1:
		sep := v[2]
		if sep != ':' && sep != '-' {
			return "", false
		}
		for k := 0; k < len(v); k++ {
			if k%3 == 2 {
				if v[k] != sep {
					return "", false
				}
				continue
			}
			c, valid := upperHex(v[k])
			if !valid {
				return "", false
			}
			out = append(out, c)
		}
	default:
		return "", false
	}
	return string(out), true
}

// upperHex returns c as an upper-case hex digit, or valid=false.
func upperHex(c byte) (byte, bool) {
	switch {
	case c >= '0' && c <= '9', c >= 'A' && c <= 'F':
		return c, true
	case c >= 'a' && c <= 'f':
		return c - 'a' + 'A', true
	default:
		return 0, false
	}
}

// hexValue returns the value of an upper-case hex digit.
func hexValue(c byte) byte {
	if c >= 'A' {
		return c - 'A' + 10
	}
	return c - '0'
}
