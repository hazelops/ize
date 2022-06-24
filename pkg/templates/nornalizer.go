package templates

import (
	"strings"
)

const Indentation = `  `

// LongDesc normalizes a command's long description to follow the conventions.
func LongDesc(s string) string {
	if len(s) == 0 {
		return s
	}
	return normalizer{s}.trim().indent().string
}

// Examples normalizes a command's examples to follow the conventions.
func Examples(s string) string {
	if len(s) == 0 {
		return s
	}
	return normalizer{s}.trim().indent().string
}

// Dedent removes any common leading whitespace from every line in text.
func Dedent(s string) string {
	if len(s) == 0 {
		return s
	}
	return normalizer{s}.trim().dedent().string
}

type normalizer struct {
	string
}

func (s normalizer) trim() normalizer {
	s.string = strings.TrimSpace(s.string)
	return s
}

func (s normalizer) dedent() normalizer {
	text := []string{}
	for _, line := range strings.Split(s.string, "\n") {
		trimmed := strings.TrimSpace(line)
		text = append(text, trimmed)
	}
	s.string = strings.Join(text, "\n")
	return s
}

func (s normalizer) indent() normalizer {
	indentedLines := []string{}
	for _, line := range strings.Split(s.string, "\n") {
		trimmed := strings.TrimSpace(line)
		indented := Indentation + trimmed
		indentedLines = append(indentedLines, indented)
	}
	s.string = strings.Join(indentedLines, "\n")
	return s
}
