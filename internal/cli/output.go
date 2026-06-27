package cli

import (
	"path/filepath"
	"strings"
)

// resolveOutput returns the directory and file path for generated code.
// If out ends with ".go", it is treated as a file path; otherwise out is a
// directory and the file is named "{pkg}.go".
func resolveOutput(out, pkg string) (dir, file string) {
	if strings.HasSuffix(out, ".go") {
		return filepath.Dir(out), out
	}
	return out, filepath.Join(out, pkg+".go")
}
