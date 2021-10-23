package env

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
)

type Env struct{}

func NewEnv() *Env {
	return &Env{}
}

// ToENV convert a map (JSON structure) to string
// (concatenated by \n)
// This function only support convert simple format
// i.e: string, number
// another types will be ignored
// it will be converted to empty string.
// The key and value will be concatenated by equal sign.
func (e *Env) ToENV(src map[string]interface{}) []byte {
	var lines []string

	for key, value := range src {
		var raw string

		switch t := value.(type) {
		case int:
			raw = addQuote(fmt.Sprintf("%v", value))
		case float64:
			raw = addQuote(fmt.Sprintf("%v", value))
		case string:
			raw = addQuote(fmt.Sprintf("%v", value))
		case bool:
			raw = addQuote(fmt.Sprintf("%v", value))
		default:
			raw = ""
			fmt.Println("WARN: This type", t, "will be ignored")
		}
		lines = append(lines, fmt.Sprintf("export %s=%s", key, raw))
	}

	content := strings.Join(lines, "\n")
	return []byte(content)
}

// addQuote to input string if it contains special characters
func addQuote(s string) string {
	var isValid = regexp.MustCompile(`^[a-zA-Z0-9.,_-]+$`).MatchString
	if !isValid(s) {
		s = strconv.Quote(s)
	}
	return s
}

// ToJSON convert string lines to a map (JSON structure)
// Each line should be formatted like export key=value or key=value
// The value can be empty. The key should contain only
// alphabet, number and underscore.
// Comment lines will be ignored.
func (e *Env) ToJSON(src []string) (map[string]interface{}, error) {
	content := make(map[string]interface{})
	var isValid = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	for _, line := range src {
		line := strings.TrimSpace(line)
		// Skip with line start with # or empty
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// Remove export
		line = strings.TrimSpace(strings.TrimPrefix(line, "export"))

		v := strings.Split(line, "=")
		// Validate key
		if !isValid(v[0]) {
			return nil, errors.New(fmt.Sprintf("Env: Key %s is invalid format", v[0]))
		}

		value := ""
		if len(v) > 1 {
			value = strings.Join(v[1:], "=")
		}

		content[v[0]] = trimQuotes(value)
	}

	return content, nil
}

// trimQuotes remove quote from string
func trimQuotes(s string) string {
	if len(s) >= 2 {
		if c := s[len(s)-1]; s[0] == c && (c == '"' || c == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
