//go:build js && wasm

package settings

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

// Path is the key of the item `name` in the page's local storage,
// under gohar's prefix.
func Path(name string) string { return "gohar/" + name }

// Load reads the item at `path` into `v`, and says whether there was
// one, as the desktop's Load does with a file. A browser that refuses
// local storage, a private window in some of them, gives no item: the
// game starts as on its first run.
func Load(path string, v any) (bool, error) {
	store := js.Global().Get("localStorage")
	if !store.Truthy() {
		return false, nil
	}
	item := store.Call("getItem", path)
	if item.IsNull() || item.IsUndefined() {
		return false, nil
	}
	return true, json.Unmarshal([]byte(item.String()), v)
}

// Save writes `v` to the item at `path` as JSON. Local storage writes
// an item whole, or not at all: there is no temporary to go through.
func Save(path string, v any) (err error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	store := js.Global().Get("localStorage")
	if !store.Truthy() {
		return nil
	}
	defer func() { // a full or refused storage throws
		if r := recover(); r != nil {
			err = fmt.Errorf("settings: %s: %v", path, r)
		}
	}()
	store.Call("setItem", path, string(data))
	return nil
}
