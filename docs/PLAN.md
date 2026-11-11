---
PLAN: "feat: input.ImageFile, DocumentFile, ScanFile, VideoFile and AudioFile — one named constructor per business asset, one numeric argument"
TAG: v0.2.0
EXECUTOR: jules
REVIEWER: none
---

> Read AGENTS.md at the repo root before any change.
> **Finish by yourself:** when the delivery check passes, commit, push and open the pull request. Never ask whether to
> proceed or finalize — nobody answers, and a session that waits ends without your work. Commit after every stage.
> Never run `git reset`, `git restore`, `git checkout -- <file>` or `git stash` on your own work.
> **Publishing:** you publish only by OPENING THE PULL REQUEST. Submit after Stage 1 and keep committing.
> **Never ask a yes/no question about progress.**
> **Before you submit:** `git status` must be clean and `git diff --stat main` must list the files of the stages.

# Plan: named file constructors in `input`

## Why (read once)

A table with a file declares its field with one expression in `definition.go`, e.g. `{Name: "logo", Type: input.ImageFile(2097152)}`. `definitionc` copies that expression into generated code and compiles a probe of it, and it cannot import the packages of the arguments (`filetype.Images`) nor read an expression like `2<<20`: so the declaration takes **one plain number** and the allowed types are chosen by the constructor's name (decision P4 = A of the file-upload design). `input.File(allowed, maxBytes)` (v0.1.0) stays for code and tests.

## Ecosystem search

`input.File` (`file.go`) and `filetype` (`Images`, `Documents`, `Scans`, `Videos`, `Audio`, published v0.1.1) exist: this plan only names combinations. No new logic.

## Design gate

1. **Prior art:** `<input type="file" accept="image/*">` names the family, not a MIME list; Rails Active Storage `has_one_attached :logo` with content-type validators by family. https://developer.mozilla.org/docs/Web/HTML/Attributes/accept (unverified in this session: from memory).
2. **Novice name:** `ImageFile`, `DocumentFile` (PDF, Word, Excel, PowerPoint), `ScanFile` (images and PDF), `VideoFile`, `AudioFile`.
3. **Ledger:** concepts +5; files to declare a file field −N; ways to do the same: +0 (`File` is the general one, the five are its named uses).
4. **Where:** `input` (the kinds live here).
5. **Deletes:** nothing.

## Tools

```bash
go install webtyp.com/devflow/cmd/gotest@latest
go install webtyp.com/lang/cmd/langc@latest
go get webtyp.com/filetype@v0.1.1
go mod tidy
```
No `replace`. If a module cannot be fetched, stop and say so in the PR.

## APIs used

| symbol | exact signature | tag | file |
|---|---|---|---|
| `input.File` | `func File(allowed filetype.Allowlist, maxBytes int) FileKind` | this repo (v0.1.0) | `file.go` |
| `filetype.Images`, `Documents`, `Scans`, `Videos`, `Audio` | `var X = NewAllowlist(...)` of type `filetype.Allowlist` | filetype v0.1.1 | `filetype.go` |
| `(FileKind).ValidateFile` | `ValidateFile(data []byte) (filetype.Type, error)` | this repo | `file.go` |

## Stage 1 — `file.go`: append these five functions

```go
// ImageFile is a file field that accepts raster images (PNG, JPEG, GIF, WebP) up to maxBytes.
//
//	{Name: "logo", Type: input.ImageFile(2097152)}
func ImageFile(maxBytes int) FileKind { return File(filetype.Images, maxBytes) }

// DocumentFile is a file field that accepts PDF, Word, Excel and PowerPoint documents up to maxBytes.
//
//	{Name: "contract", Type: input.DocumentFile(10485760)}
func DocumentFile(maxBytes int) FileKind { return File(filetype.Documents, maxBytes) }

// ScanFile is a file field that accepts scanned papers and photographed receipts: raster images and PDF, up to maxBytes.
//
//	{Name: "receipt", Type: input.ScanFile(5242880)}
func ScanFile(maxBytes int) FileKind { return File(filetype.Scans, maxBytes) }

// VideoFile is a file field that accepts short videos (MP4, MOV, WebM) up to maxBytes. The whole body is read in memory
// on arrival: keep it small (a few tens of megabytes), never a long recording.
//
//	{Name: "clip", Type: input.VideoFile(26214400)}
func VideoFile(maxBytes int) FileKind { return File(filetype.Videos, maxBytes) }

// AudioFile is a file field that accepts audio (MP3, WAV, OGG, M4A) up to maxBytes.
//
//	{Name: "note", Type: input.AudioFile(10485760)}
func AudioFile(maxBytes int) FileKind { return File(filetype.Audio, maxBytes) }
```

## Stage 2 — tests (`tests/file_kinds_test.go`, package `tests`; one row each in `docs/TESTS.md`, `| [Name](../tests/file_kinds_test.go) | rule |`)

Sample bytes (each ends with two zero bytes): `png := []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A, 0, 0}`; `pdf := []byte("%PDF-1.7")`; `mp4 := []byte{0, 0, 0, 0x18, 'f', 't', 'y', 'p', 'i', 's', 'o', 'm', 0, 0}`; `wav := []byte("RIFF\x24\x00\x00\x00WAVEfmt ")`.

| Test | setup → call → expected |
|---|---|
| `TestImageFile_AcceptsImageRejectsPDF` | `input.ImageFile(1024).ValidateFile(png)` ok; `ValidateFile(pdf)` → error |
| `TestDocumentFile_AcceptsPDFRejectsImage` | `input.DocumentFile(1024).ValidateFile(pdf)` ok; `ValidateFile(png)` → error |
| `TestScanFile_AcceptsImageAndPDF` | `input.ScanFile(1024)`: `ValidateFile(png)` ok and `ValidateFile(pdf)` ok |
| `TestVideoFile_AcceptsMP4RejectsAudio` | `input.VideoFile(1024).ValidateFile(mp4)` ok; `ValidateFile(wav)` → error |
| `TestAudioFile_AcceptsWAVRejectsVideo` | `input.AudioFile(1024).ValidateFile(wav)` ok; `ValidateFile(mp4)` → error |
| `TestFileKinds_StoreText` | each of the five: `Storage()` is `model.FieldText`, `Name()` is `"file"`, `MaxBytes()` is the number given |

## Stage 3 — documentation (English, `webtyp/*`)

`README.md` (≤150 lines): the five constructors in the list of kinds, one line each with the example of its doc comment. No new `fmt.Err` text: `langc sync .` must report nothing new.

## Delivery check

`gotest` passes on every destination; `langc sync .` is clean; `docs/TESTS.md` has the 6 rows; the pull request is open. A tree that does not compile or whose tests did not run is a failed delivery.
