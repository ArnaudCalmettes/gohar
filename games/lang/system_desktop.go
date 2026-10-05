//go:build !js

package lang

import (
	"os"
	"strings"
)

// System returns the language of the user's session, from the usual
// variables, the fallback when none says.
func System() string {
	for _, v := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		s := os.Getenv(v)
		if s == "" || s == "C" || s == "POSIX" {
			continue
		}
		s, _, _ = strings.Cut(s, ".") // fr_FR.UTF-8
		return strings.ReplaceAll(s, "_", "-")
	}
	return Fallback
}
