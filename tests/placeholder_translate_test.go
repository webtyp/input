//go:build wasm

package input_test

import (
	"os"
	"syscall/js"
	"testing"

	"webtyp.com/input"
	"webtyp.com/lang"
)

// The page bootstrap preloads the dictionary into globalThis.__webtypLang.
func TestMain(m *testing.M) {
	parsed := js.Global().Get("JSON").Call("parse", `{"default":"es","languages":["es"],"keys":{"example:":["ejemplo:"]}}`)
	js.Global().Set(lang.GlobalDictKey, parsed)
	os.Exit(m.Run())
}

// GetPlaceholder translates each part on every call, so the active language
// applies without building a new input.
func TestIPPlaceholderTranslated(t *testing.T) {
	getter, ok := input.IP().(interface{ GetPlaceholder() string })
	if !ok {
		t.Fatal("IP() must expose GetPlaceholder")
	}
	lang.OutLang(lang.ES)
	defer lang.OutLang(lang.EN)
	if got := getter.GetPlaceholder(); got != "ejemplo: 192.168.1.1" {
		t.Errorf("got %q, want %q", got, "ejemplo: 192.168.1.1")
	}
}
