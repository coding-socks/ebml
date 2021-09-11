package ebmlpath

import (
	"strings"
)

// Join joins any number of path elements into a single path,
// separating them with backslashes. Empty elements are ignored.
func Join(elem ...string) string {
	size := 0
	for _, e := range elem {
		size += len(e)
	}
	if size == 0 {
		return ""
	}
	buf := make([]byte, 0, size+len(elem)-1)
	for _, e := range elem {
		if len(buf) > 0 || e != "" {
			if len(buf) > 0 || !strings.HasPrefix(e, "\\") {
				buf = append(buf, '\\')
			}
			buf = append(buf, e...)
		}
	}
	return string(buf)
}
