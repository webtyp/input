# Tests

| Test | setup → call → expected |
|---|---|
| [TestFile_AcceptsAllowedTypeFromBytes](../tests/file_test.go) | `k := input.File(filetype.Images, 1024)`; `k.ValidateFile(png)` → `filetype.PNG`, no error |
| [TestFile_IgnoresWhatTheClientClaims](../tests/file_test.go) | the SVG bytes with `k := input.File(filetype.NewAllowlist(filetype.PNG), 1024)` → error (the type comes from the bytes) |
| [TestFile_RejectsSVG](../tests/file_test.go) | `k.ValidateFile(svgBytes)` with `filetype.Images` → error whose text contains `SVG` |
| [TestFile_RejectsOverMaxBytes](../tests/file_test.go) | `input.File(filetype.Images, 4).ValidateFile(png)` (10 bytes) → error; `MaxBytes()` is 4 |
| [TestFile_KeyValidation](../tests/file_test.go) | `Validate("01HX9-a_b.png")` → nil; `Validate("../x.png")` → error; `Validate("a/b.png")` → error; `Validate("")` with the kind required → error |
| [TestFile_StorageIsText](../tests/file_test.go) | `k.Storage()` is `model.FieldText`, `k.Name()` is `"file"` |
| [TestFile_CloneKeepsTheRules](../tests/file_test.go) | `c := k.Clone("form", "logo").(input.FileKind)`; `c.MaxBytes()` is the same; `c.ValidateFile(svgBytes)` still errors |
| [TestCanonicalIP](../inputs_test.go) | |
| [TestClone_Preservation](../inputs_test.go) | |
| [TestIP_Placeholder](../inputs_test.go) | |
| [TestMasked_DefaultsFalse](../mask_test.go) | |
| [TestMasked_DoesNotAffectValidation](../mask_test.go) | |
| [TestMasked_SurvivesClone](../mask_test.go) | |
| [TestPassword_IsMaskedByDefault](../mask_test.go) | |
| [TestCanonicalMAC](../tests/mac_test.go) | |
| [TestIsLocallyAdministeredMAC](../tests/mac_test.go) | |
| [TestMACClonedField](../tests/mac_test.go) | |
| [TestMACValidateAccepts](../tests/mac_test.go) | |
| [TestMACValidateEmptyMatchesIP](../tests/mac_test.go) | |
| [TestMACValidateRejects](../tests/mac_test.go) | |
| [TestIPPlaceholderTranslated](../tests/placeholder_translate_test.go) | |
| [TestNormalizeRUT_AnyWayItIsWritten](../tests/rut_normalize_test.go) | |
| [TestSelectWithOptionsChecksMembership](../tests/select_test.go) | |
| [TestSelectWithoutOptionsAcceptsIDs](../tests/select_test.go) | |
| [Test_Validation](../validation_test.go) | |
