//go:build wasm

package input_test

import (
	"os"
	"syscall/js"
	"testing"

	"webtyp.com/input"
	"webtyp.com/lang"
)

// The page carries the dictionary (sitec inlines it); insert it before any lookup.
func TestMain(m *testing.M) {
	doc := js.Global().Get("document")
	el := doc.Call("createElement", "script")
	el.Set("type", "application/json")
	el.Set("id", lang.ScriptID)
	el.Set("textContent", `{"default":"es","languages":["es"],"keys":{"example:":["ejemplo:"]}}`)
	doc.Get("head").Call("appendChild", el)
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
