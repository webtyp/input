package input

import "webtyp.com/fmt"

type ip struct{ Base }

// IP creates a new IP input instance.
func IP() Input {
	i := &ip{}
	i.Numbers = true
	i.Letters = true // hex for ipv6
	i.Extra = []rune{'.', ':'}
	i.Minimum = 7  // 1.1.1.1
	i.Maximum = 39 // full ipv6 length
	i.InitBase("", "", "text")
	// The type accepts IPv4 AND IPv6, so a user cannot infer the shape from
	// the field name — the format hint belongs to the type that defines the
	// format, the precedent rut.go already set.
	i.SetPlaceholder("example:", "192.168.1.1")
	return i
}

// Validate validates IPv4 or IPv6 format.
func (i *ip) Validate(value string) error {
	if value == "0.0.0.0" {
		return fmt.Err("Format", "Invalid")
	}

	origMin := i.Minimum
	if value != "" {
		hasColon := false
		for _, c := range value {
			if c == ':' {
				hasColon = true
				break
			}
		}
		if hasColon {
			i.Minimum = 2
		}
	}

	err := i.Permitted.Validate(i.name, value)
	i.Minimum = origMin
	if err != nil {
		return err
	}

	hasDot, hasColon := false, false
	for _, c := range value {
		if c == '.' {
			hasDot = true
		}
		if c == ':' {
			hasColon = true
		}
	}

	switch {
	case hasDot && hasColon:
		return fmt.Err("Format", "Invalid")
	case !hasDot && !hasColon:
		return fmt.Err("Format", "Invalid")
	case hasDot:
		return validateIPv4(value)
	default:
		return validateIPv6(value)
	}
}

func validateIPv4(value string) error {
	parts := fmt.Convert(value).Split(".")
	if len(parts) != 4 {
		return fmt.Err("Format", "Invalid")
	}
	for _, part := range parts {
		if len(part) < 1 || len(part) > 3 {
			return fmt.Err("Format", "Invalid")
		}
		for _, c := range part {
			if c < '0' || c > '9' {
				return fmt.Err("Format", "Invalid")
			}
		}
		val, err := fmt.Convert(part).Int()
		if err != nil || val < 0 || val > 255 {
			return fmt.Err("Format", "Invalid")
		}
	}
	return nil
}

func validateIPv6(value string) error {
	// Count occurrences of "::"
	doubleColonCount := 0
	for i := 0; i < len(value)-1; i++ {
		if value[i] == ':' && value[i+1] == ':' {
			doubleColonCount++
			i++
		}
	}

	if doubleColonCount > 1 {
		return fmt.Err("Format", "Invalid")
	}

	parts := fmt.Convert(value).Split(":")

	startsWithDoubleColon := len(value) >= 2 && value[0] == ':' && value[1] == ':'
	endsWithDoubleColon := len(value) >= 2 && value[len(value)-2] == ':' && value[len(value)-1] == ':'

	if doubleColonCount == 1 {
		if startsWithDoubleColon && endsWithDoubleColon {
			// e.g. "::" -> parts is ["", "", ""]
			if len(parts) != 3 {
				return fmt.Err("Format", "Invalid")
			}
			for _, part := range parts {
				if part != "" {
					return fmt.Err("Format", "Invalid")
				}
			}
		} else if startsWithDoubleColon {
			// e.g. "::1" -> parts is ["", "", "1"]
			if len(parts) < 3 {
				return fmt.Err("Format", "Invalid")
			}
			if parts[0] != "" || parts[1] != "" {
				return fmt.Err("Format", "Invalid")
			}
			nonEmptyCount := 0
			for i := 2; i < len(parts); i++ {
				if parts[i] == "" {
					return fmt.Err("Format", "Invalid")
				}
				if !isValidHexPart(parts[i]) {
					return fmt.Err("Format", "Invalid")
				}
				nonEmptyCount++
			}
			if nonEmptyCount > 7 {
				return fmt.Err("Format", "Invalid")
			}
		} else if endsWithDoubleColon {
			// e.g. "1::" -> parts is ["1", "", ""]
			if len(parts) < 3 {
				return fmt.Err("Format", "Invalid")
			}
			if parts[len(parts)-1] != "" || parts[len(parts)-2] != "" {
				return fmt.Err("Format", "Invalid")
			}
			nonEmptyCount := 0
			for i := 0; i < len(parts)-2; i++ {
				if parts[i] == "" {
					return fmt.Err("Format", "Invalid")
				}
				if !isValidHexPart(parts[i]) {
					return fmt.Err("Format", "Invalid")
				}
				nonEmptyCount++
			}
			if nonEmptyCount > 7 {
				return fmt.Err("Format", "Invalid")
			}
		} else {
			// e.g. "2001::1" -> parts has exactly one empty part
			emptyCount := 0
			nonEmptyCount := 0
			for _, part := range parts {
				if part == "" {
					emptyCount++
				} else {
					if !isValidHexPart(part) {
						return fmt.Err("Format", "Invalid")
					}
					nonEmptyCount++
				}
			}
			if emptyCount != 1 || nonEmptyCount > 7 {
				return fmt.Err("Format", "Invalid")
			}
		}
	} else {
		// doubleColonCount == 0
		if len(parts) != 8 {
			return fmt.Err("Format", "Invalid")
		}
		for _, part := range parts {
			if part == "" || !isValidHexPart(part) {
				return fmt.Err("Format", "Invalid")
			}
		}
	}

	return nil
}

func isValidHexPart(part string) bool {
	if len(part) < 1 || len(part) > 4 {
		return false
	}
	for _, c := range part {
		isDigit := c >= '0' && c <= '9'
		isLowerHex := c >= 'a' && c <= 'f'
		isUpperHex := c >= 'A' && c <= 'F'
		if !isDigit && !isLowerHex && !isUpperHex {
			return false
		}
	}
	return true
}

// Clone satisfies input.Input — IP() returns Input which implements it.
func (i *ip) Clone(parentID, name string) Input {
	c := *i
	c.InitBase(parentID, name, "text")
	return &c
}

// loopbackIP is the one spelling CanonicalIP gives every loopback address.
// IPv4, because that is what an operator types into a device form and what a
// log reader recognises; the IPv6 loopback carries no extra information.
const loopbackIP = "127.0.0.1"

// ipv4MappedPrefix is how a dual-stack socket reports an IPv4 peer
// ("::ffff:192.168.1.10"); the address is the IPv4 part.
const ipv4MappedPrefix = "::ffff:"

// CanonicalIP returns the single spelling of an IP address that two IPs are
// compared by. Whoever stores an IP and whoever later looks it up must both
// pass it through here, or the same machine is a different string depending on
// how it connected: one "localhost" reaches a server as ::1 from one client and
// as 127.0.0.1 from another.
//
//   - surrounding space is trimmed and hex digits are lowercased;
//   - an IPv4-mapped IPv6 address becomes its IPv4 address;
//   - every loopback address (::1 in any spelling, 127.0.0.0/8) becomes
//     127.0.0.1 — they all name this same machine, so merging them grants no
//     other host anything.
//
// Any other value is returned trimmed and lowercased, not validated: IP().Validate
// is the validator. Zero-compression of other IPv6 addresses is not rewritten.
func CanonicalIP(value string) string {
	v := fmt.Convert(value).TrimSpace().ToLower().String()
	if len(v) > len(ipv4MappedPrefix) && v[:len(ipv4MappedPrefix)] == ipv4MappedPrefix {
		if mapped := v[len(ipv4MappedPrefix):]; validateIPv4(mapped) == nil {
			v = mapped
		}
	}
	if isLoopbackIPv4(v) || isLoopbackIPv6(v) {
		return loopbackIP
	}
	return v
}

func isLoopbackIPv4(v string) bool {
	return len(v) > 4 && v[:4] == "127." && validateIPv4(v) == nil
}

// isLoopbackIPv6 reports whether v is ::1 in any spelling: every group before
// the last is zero (or elided by "::") and the last group's value is 1.
func isLoopbackIPv6(v string) bool {
	last := -1
	for i := len(v) - 1; i >= 0; i-- {
		if v[i] == ':' {
			last = i
			break
		}
	}
	if last == -1 || validateIPv6(v) != nil {
		return false
	}
	for _, c := range v[:last] {
		if c != '0' && c != ':' {
			return false
		}
	}
	group := v[last+1:]
	for len(group) > 1 && group[0] == '0' {
		group = group[1:]
	}
	return group == "1"
}
