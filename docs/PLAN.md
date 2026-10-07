---
PLAN: "feat: MAC input kind, CanonicalMAC and IsLocallyAdministeredMAC"
EXECUTOR: jules
REVIEWER: none
STATUS: running
SESSION: 8409417972468080869
---

> This plan is dispatched via the CodeJob workflow. See skill: agents-workflow.
>
> Phase F1 of the network administration master plan
> (`veltylabs/mjosefa-cms` → `docs/RED_ADMINISTRACION_MASTER_PLAN.md`, private repo — you do not
> need it; everything required is in this file).

# Plan — `webtyp.com/input`: a MAC address input

## Development rules (apply to every line you write)

- **Only `webtyp.com/fmt`.** No `errors`, `strconv`, `strings`, stdlib `fmt`. This package is
  compiled to WASM/TinyGo.
- **No Go `map`**, no `reflect`.
- **No string literals in logic**: every repeated string is a named constant.
- **Tests go in `tests/`** (package `input_test`, external, public API only). Do not add root-level
  tests. Never export a symbol only so a test can reach it.
- Runner is `gotest` (never `go test`). Both `gotest ./...` and the WASM suite must be green.

## 0. Context

A downstream inventory (`github.com/veltylabs/device_manager`) is starting to identify every network
interface by its **MAC address** (the 6-byte hardware id of a network card), because an IP address can
change and a MAC is what a router's DHCP server and firewall key on. The form for that field needs a
`Kind` with a widget, exactly like `input.IP()` already provides for IPs (`ip.go` — read it, this plan
mirrors it).

Two facts about MAC addresses drive the API:

1. **The same MAC is written three ways** depending on who prints it: `48:F1:7F:D9:D7:B7` (RouterOS,
   Linux), `48-F1-7F-D9-D7-B7` (Windows `ipconfig /all`), `48F17FD9D7B7` (PowerShell
   `Get-NetAdapter … PermanentAddress`), in upper or lower case. Whoever stores a MAC and whoever looks
   it up must agree on one spelling — the same reason `CanonicalIP` exists.
2. **Phones and Windows can invent a MAC per Wi-Fi network** ("private address", "random hardware
   addresses"; Windows can even change it daily). Such a MAC has the *locally administered* bit set:
   bit `0x02` of the first byte, i.e. the second hex digit is `2`, `6`, `A` or `E`. A real incident
   caused by this is why the consumer must be able to detect it: a PC registered with its real MAC
   came back with an invented one and lost access.

## Design gate

**1. Prior art.**
- Go stdlib `net.ParseMAC` accepts `:`/`-`/`.` separated forms and returns a `HardwareAddr` whose
  `String()` is lowercase colon form — one canonical spelling, the model followed here (but stdlib is
  banned in this package).
- HTML has no `type="mac"`; form libraries (e.g. Django's `django-macaddress`, the `validator.js`
  `isMACAddress`) treat it as a text input with a format validator plus a normalizer — exactly the
  `input.IP()` + `CanonicalIP` pair this package already has.
- IEEE 802 defines the locally administered bit; Android/iOS/Windows MAC randomization sets it, and
  network tools (Wireshark, `macchanger`) report it under that name.
This package differs only in that it may not use stdlib, so parsing is hand-written over `webtyp/fmt`.

**2. Novice-name test.** `input.MAC()` reads like `input.IP()` next to it. `input.CanonicalMAC(v)`
mirrors `input.CanonicalIP(v)`. `input.IsLocallyAdministeredMAC(v)` uses the IEEE term that tools
print; its doc comment says "true for randomized/private addresses" so a reader who knows only that
word finds it.

**3. Complexity ledger.**
```
Concepts the developer must learn   +1 (MAC kind; its helpers mirror the IP ones)
Files they must touch to do X       0
Lines at the call site              {Type: input.MAC()} — 1
Ways to do the same thing           0 (no MAC handling exists anywhere today)
```

**4. Where it belongs.** A format validator + widget for a form field is exactly this package's
concern (`IP`, `Rut`, `Email` live here). The *policy* "reject randomized MACs" is the consumer's
(device_manager decides); this package only answers the factual question.

**5. What it deletes.** Nothing — new capability.

## 1. Stage 1 — `mac.go` (new file, package `input`)

```go
package input

import "webtyp.com/fmt"

type mac struct{ Base }

// macSeparator is the separator CanonicalMAC writes.
const macSeparator = ':'

// macHexDigits is the number of hex digits in a MAC address (6 bytes).
const macHexDigits = 12

// MAC creates a MAC address input. It accepts "AA:BB:CC:DD:EE:FF",
// "AA-BB-CC-DD-EE-FF" and "AABBCCDDEEFF", in any case.
func MAC() Input {
	i := &mac{}
	i.Numbers = true
	i.Letters = true // hex digits a-f
	i.Extra = []rune{':', '-'}
	i.Minimum = macHexDigits     // "AABBCCDDEEFF"
	i.Maximum = macHexDigits + 5 // "AA:BB:CC:DD:EE:FF"
	i.InitBase("", "", "text")
	i.SetPlaceholder("example:", "48:F1:7F:D9:D7:B7")
	return i
}
```

Implement on `*mac`:

- `Validate(value string) error` — first `i.Permitted.Validate(i.name, value)` (same as `ip.go`), then
  `parseMAC(value)`; any failure returns `fmt.Err("Format", "Invalid")` (the exact error `ip.go` uses).
- `Clone(parentID, name string) Input` — copy of `ip.go`'s `Clone`, with `"text"`.

Unexported helper (single parser, used by all three exported functions — never duplicate it):

```go
// parseMAC returns the 12 hex digits of value in upper case, or ok=false.
// Accepted shapes: 12 hex digits; or 6 groups of 2 hex digits separated by
// ':' everywhere or by '-' everywhere (mixed separators are rejected).
func parseMAC(value string) (digits string, ok bool)
```
Trim surrounding space first (`fmt.Convert(value).TrimSpace().String()`).

Exported functions (same file):

```go
// CanonicalMAC returns the single spelling two MACs are compared by:
// upper case, colon-separated ("48:F1:7F:D9:D7:B7"). Whoever stores a MAC
// and whoever later looks it up must both pass it through here.
// A value that is not a MAC is returned trimmed and upper-cased, unchanged
// otherwise — MAC().Validate is the validator.
func CanonicalMAC(value string) string

// IsLocallyAdministeredMAC reports whether the MAC has the IEEE "locally
// administered" bit set (second hex digit 2, 6, A or E). That is the case
// for randomized/private addresses that phones and Windows invent per
// Wi-Fi network, and for virtual machines. False for an invalid MAC.
func IsLocallyAdministeredMAC(value string) bool
```

## 2. Stage 2 — tests (`tests/mac_test.go`, no build tag, package `input_test`)

Table tests, asserting results (never discard returns):

| Call | Input | Expected |
|---|---|---|
| `MAC().Validate` | `48:F1:7F:D9:D7:B7`, `48-f1-7f-d9-d7-b7`, `48F17FD9D7B7`, ` 48:f1:7f:d9:d7:b7 ` | `nil` |
| `MAC().Validate` | `48:F1:7F:D9:D7`, `48:F1:7F:D9:D7:B7:00`, `48:F1-7F:D9:D7:B7` (mixed), `GG:F1:7F:D9:D7:B7`, `48F17FD9D7B` | non-nil |
| `CanonicalMAC` | each valid input above | `48:F1:7F:D9:D7:B7` |
| `IsLocallyAdministeredMAC` | `4A:C5:93:7A:12:DE`, `DE:94:42:52:C5:59`, `12:AB:2F:3F:6C:A3`, `CE:45:90:5F:05:E7` | `true` |
| `IsLocallyAdministeredMAC` | `48:F1:7F:D9:D7:B7`, `60:6C:66:C0:3F:1E`, `not-a-mac` | `false` |

Note on `""`: follow exactly what `IP().Validate("")` does today (read `ip.go`: `Permitted.Validate`
with `Minimum` set). Write the test to the same behaviour and say so in a comment.

Add one consumer-shaped case: a `model.Definition`-style field `{Name: "mac", Type: input.MAC()}`
cloned with `Clone("form1", "mac")` validates `48-F1-7F-D9-D7-B7` and rejects `48:F1`.

## 3. Stage 3 — docs

- `README.md` → "Available Inputs" table: add the row
  `| \`MAC\` | \`text\` | 12 hex digits, plain or \`:\`/\`-\` separated; see \`CanonicalMAC\`, \`IsLocallyAdministeredMAC\` |`
  in alphabetical position.

## Acceptance criteria

- `gotest ./...` green (including the WASM run).
- `grep -rn '"strings"\|"strconv"\|"errors"' mac.go` → empty.
- `grep -rn "map\[" mac.go tests/mac_test.go` → empty.

| Stage | Files | Done when |
|---|---|---|
| 1 | `mac.go` | compiles, no stdlib |
| 2 | `tests/mac_test.go` | table above green |
| 3 | `README.md` | row added |
