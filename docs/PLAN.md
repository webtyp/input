---
PLAN: "feat: input.File — the kind of a file field, validated by its bytes, the same in the browser and in crud"
TAG: v0.1.0
EXECUTOR: jules
REVIEWER: none
STATUS: running
SESSION: 947645065135672486
---

> Read AGENTS.md at the repo root before any change.
> **Finish by yourself:** when the delivery check passes, commit, push and open the pull request. Never ask whether to
> proceed or finalize — nobody answers, and a session that waits ends without your work. Commit after every stage.
> Never run `git reset`, `git restore`, `git checkout -- <file>` or `git stash` on your own work.
> **Publishing:** you publish only by OPENING THE PULL REQUEST (a `git push` from your sandbox does not reach GitHub). Open it after Stage 2 with whatever compiles and keep committing to it.
> **Never ask a yes/no question about progress**; ask only when a signature in this plan is wrong, and name the symbol.
> **Before you submit:** `git status` must be clean and `git diff --stat main` must list the files of the stages. A pull request that carries only go.mod, go.sum or docs/PLAN.md is a failed delivery: `git add -A`, commit and submit again.

# Plan: `input.File`

## Why (read once)

A table with a file (the logo of a site, the photo of a patient) stores only the file's KEY in a text column; the bytes travel in the body of the request (design D46). Who validates the bytes? The kind of the field, so the same code runs in the browser before sending and in `crud` when the request arrives: a type read from the bytes (`filetype`), never from the client, and a maximum size. This plan adds that kind. Nothing else (no storage, no route, no widget).

## Ecosystem search

Searched `filetype`, `Allowlist`, `Detect`, `MaxSize` over `webtyp` and `veltylabs`: `filetype` (`filetype.Detect`, `filetype.Allowlist.Validate`, `filetype.Images`) already reads the type from the bytes and rejects SVG/HTML by name; `cloudflare/files` has its own size check. **Reuse** `filetype`; **new** here only the kind that carries the allowlist and the size in a field declaration. `input.Select` / `input.Email` (`select.go`, `email.go`) are the shape to copy.

## Design gate

1. **Prior art:** HTML `<input type="file" accept="…">` (the browser filters by declared type, which is a claim); Rails Active Storage `content_type` validators and Django `FileExtensionValidator` (both validate the name or the declared type, unlike this one); Stripe file uploads (type checked on the server from the content). https://developer.mozilla.org/docs/Web/HTML/Element/input/file (unverified in this session: from memory). We differ on purpose: the type comes from the bytes (D46).
2. **Novice name:** `input.File(allowed, maxBytes)`: the allowed kinds and the largest size, read aloud as a sentence.
3. **Ledger:** concepts +1 (`File`), files to touch to add a validated file field −N (one line in `definition.go`), lines at the call site +0, ways to do the same: 0.
4. **Where:** `input` (the kinds live here); `filetype` is the only new import; no `internal/`.
5. **Deletes:** nothing (new capability).

## Tools

```bash
go install webtyp.com/devflow/cmd/gotest@latest
go install webtyp.com/lang/cmd/langc@latest
go get webtyp.com/model@v0.2.20 webtyp.com/filetype@latest webtyp.com/lang@latest
go mod tidy
```
No `replace` in `go.mod`. If a module cannot be fetched, stop and say so in the PR; never write a hash by hand.

## APIs used

| symbol | exact signature | tag | file |
|---|---|---|---|
| `filetype.Allowlist` | `type Allowlist struct{…}`; `func NewAllowlist(t ...Type) Allowlist` | filetype latest | `filetype.go` |
| `(Allowlist).Validate` | `func (a Allowlist) Validate(data []byte) (Type, error)` | filetype latest | `filetype.go:104` |
| `filetype.Images` | `var Images = NewAllowlist(PNG, JPEG, GIF, WebP)` | filetype latest | `filetype.go` |
| `filetype.Type` | `struct{ MIME string; Ext string; … }` with `String() string` | filetype latest | `filetype.go` |
| `input.Base` | `InitBase(parentID, name, htmlName string)`; fields `Letters`, `Numbers`, `Extra []rune`, `Minimum`, `Maximum` | this repo | `base.go` |
| `model.FieldText` | the storage of a text column | model v0.2.20 | `model/kind.go` |
If a signature in the published tag differs from this table, stop and name the symbol in the PR; never adapt it by guessing.

## Stage 1 — `file.go` (new)

```go
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
```

`Validate(value string)` is inherited from `Base`: it checks the stored KEY against the permitted characters and lengths (so `a/b` and `../x` are refused). Run `langc sync .`: the new words `File` and `Large` get their rows in `lang.json`; fill EVERY empty value (all languages the file lists) and reuse words already there (**langc**); `langc sync .` must report zero untranslated keys and zero reuse warnings.

## Stage 2 — tests (`tests/file_test.go`, external package `tests`; one row each in `docs/TESTS.md`, `| [Name](../tests/file_test.go) | rule |`)

PNG bytes for the tests: `png := []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A, 0, 0}`; SVG bytes: `[]byte("<svg xmlns='http://www.w3.org/2000/svg'></svg>")`.

| Test | setup → call → expected |
|---|---|
| `TestFile_AcceptsAllowedTypeFromBytes` | `k := input.File(filetype.Images, 1024)`; `k.ValidateFile(png)` → `filetype.PNG`, no error |
| `TestFile_IgnoresWhatTheClientClaims` | the SVG bytes with `k := input.File(filetype.NewAllowlist(filetype.PNG), 1024)` → error (the type comes from the bytes) |
| `TestFile_RejectsSVG` | `k.ValidateFile(svgBytes)` with `filetype.Images` → error whose text contains `SVG` |
| `TestFile_RejectsOverMaxBytes` | `input.File(filetype.Images, 4).ValidateFile(png)` (10 bytes) → error; `MaxBytes()` is 4 |
| `TestFile_KeyValidation` | `Validate("01HX9-a_b.png")` → nil; `Validate("../x.png")` → error; `Validate("a/b.png")` → error; `Validate("")` with the kind required → error |
| `TestFile_StorageIsText` | `k.Storage()` is `model.FieldText`, `k.Name()` is `"file"` |
| `TestFile_CloneKeepsTheRules` | `c := k.Clone("form", "logo").(input.FileKind)`; `c.MaxBytes()` is the same; `c.ValidateFile(svgBytes)` still errors |

## Stage 3 — documentation (English, `webtyp/*`)

`README.md` (≤150 lines): add `File` to the list of kinds with the one-line example above. `docs/ARCHITECTURE.md`: if it lists the kinds, add `File` and say that its bytes are validated by `ValidateFile` (browser and `crud`) and that SVG/HTML are never allowed; otherwise add one sentence to the README only. Do not create other docs.

## Delivery check

`gotest` passes on every destination; `langc sync .` reports zero untranslated keys and zero reuse warnings; `docs/TESTS.md` has the 7 rows above; the pull request is open. A tree that does not compile or whose tests did not run is a failed delivery.
