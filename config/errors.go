package config

import (
	"fmt"
	"strings"
)

type KeyError struct {
	Key string
}

func (e *KeyError) Error() string {
	return fmt.Sprintf("missing required key: %s", e.Key)
}

type PathError struct {
	Key     string
	Segment string
}

func (e *PathError) Error() string {
	return fmt.Sprintf("invalid path segment '%s' in key: %s", e.Segment, e.Key)
}

type TypeError struct {
	Key      string
	Expected string
	Actual   string
}

func (e *TypeError) Error() string {
	return fmt.Sprintf("invalid type for key '%s': expected %s, got %s", e.Key, e.Expected, e.Actual)
}

type MultiError struct {
	Errors []error
}

func (m MultiError) Error() string {
	var sb strings.Builder
	sb.WriteString("multiple errors: [")
	for i, err := range m.Errors {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(err.Error())
	}
	sb.WriteString("]")
	return sb.String()
}
