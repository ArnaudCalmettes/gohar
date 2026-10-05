//go:build js && wasm

package lang

import "syscall/js"

// System returns the language the browser is set to, "fr-FR" for
// instance, the fallback when it says none.
func System() string {
	nav := js.Global().Get("navigator")
	if !nav.Truthy() {
		return Fallback
	}
	if l := nav.Get("language"); l.Type() == js.TypeString && l.String() != "" {
		return l.String()
	}
	return Fallback
}
