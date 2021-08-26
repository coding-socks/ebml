package ebmlpath

import (
	"errors"
	"fmt"
	"strconv"
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

// PathExp is the representation of a compiled path expression.
type PathExp struct {
	expr  string
	nodes []node
}

// String returns the source text used to compile the path expression.
func (pe PathExp) String() string {
	return pe.expr
}

// Match reports whether the string s
// contains any match of the path expression pe.
func (pe PathExp) Match(s string) bool {
	if s == "" {
		return len(pe.nodes) == 0
	}
	if pe.expr == s { // low-hanging fruit
		return true
	}
	if s[0] != '\\' {
		return false
	}
	segments := strings.Split(s[1:], `\`)
	j := 0
	m := j
	for i := range pe.nodes {
		switch n := pe.nodes[i].(type) {
		case pathNode:
			found := false
			for ; j <= m && j < len(segments); j++ {
				if n.name == segments[j] {
					found = true
					break
				}
			}
			if n.recursive {
				for ; (j + 1) < len(segments); j++ {
					if segments[j] != segments[j+1] {
						break
					}
				}
			}
			if !found {
				return false
			}
			j++
			m = j
		case placeholderNode:
			j += n.minOccur
			if j >= len(segments) {
				return false
			}
			if n.maxOccur == 0 {
				m = len(segments) - 1
			} else {
				m = j + n.maxOccur
			}
		}
	}
	return j == len(segments)
}

// Match reports whether the string s
// contains any match of the path expression pattern.
func Match(pattern, s string) (bool, error) {
	c, err := Compile(pattern)
	if err != nil {
		return false, err
	}
	return c.Match(s), nil
}

const (
	stateMinOccur int = iota
	stateMaxOccur
)

// Compile parses a path expression and returns, if successful,
// a PathExp object that can be used to match against text.
func Compile(pattern string) (*PathExp, error) {
	if pattern[0] != '\\' {
		return nil, errors.New(`ebmlpath: pattern has to start with \`)
	}
	var e PathExp
	e.expr = pattern
	pattern = pattern[1:]
	for i := 0; i < len(pattern); {
		j := i
		if pattern[i] == '(' {
			j++ // skip open bracket
			var s = stateMinOccur
			var minOccur, maxOccur strings.Builder
			for ; j < len(pattern); j++ {
				if pattern[j] == ')' {
					break
				}
				if pattern[j] == '\\' {
					continue
				}
				if s == stateMinOccur && pattern[j] == '-' {
					s = stateMaxOccur
					continue
				}
				if !IsNumeric(pattern[j]) {
					return nil, fmt.Errorf("ebmlpath: unexpected non numeric value in %s at position %d", pattern, j)
				}
				switch s {
				case stateMinOccur:
					minOccur.WriteByte(pattern[j])
				case stateMaxOccur:
					maxOccur.WriteByte(pattern[j])
				}
			}
			min, _ := strconv.Atoi(minOccur.String())
			max, _ := strconv.Atoi(maxOccur.String())
			e.nodes = append(e.nodes, placeholderNode{
				minOccur: min,
				maxOccur: max,
			})
		} else {
			for ; j < len(pattern); j++ {
				if pattern[j] == '\\' {
					break
				}
			}
			if pattern[i] == '+' {
				e.nodes = append(e.nodes, pathNode{name: pattern[i+1 : j], recursive: true})
			} else {
				e.nodes = append(e.nodes, pathNode{name: pattern[i:j]})
			}
		}
		i = j + 1
	}
	if e.nodes[len(e.nodes)-1].Type() == typePlaceholder {
		return nil, errors.New("ebmlpath: pattern must not end with placeholder")
	}
	return &e, nil
}

const (
	typePath nodeType = iota + 1
	typePlaceholder
)

// A nodeType indicates what type a node belongs to.
type nodeType int

type node interface {
	// Type returns a type of this node.
	Type() nodeType
}

type pathNode struct {
	name      string
	recursive bool
}

func (p pathNode) Type() nodeType {
	return typePath
}

type placeholderNode struct {
	minOccur, maxOccur int
}

func (p placeholderNode) Type() nodeType {
	return typePlaceholder
}

// IsNumeric returns true if the given character is a numeric.
func IsNumeric(c byte) bool {
	return c >= '0' && c <= '9'
}
