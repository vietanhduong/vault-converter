package env

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

type Env struct{}

func NewEnv() *Env {
	return &Env{}
}

// ToENV convert a map (JSON structure) to string
// (concatenated by \n)
// This function only support convert simple format
// i.e: string, number, integer array, string array
// If the value is a type of map[string]interface{}
// it will be converted to empty string.
// The key and value will be concatenated by equal sign.
// If the value have array format. The elements in value
// will be concatenated by comma.
func (e *Env) ToENV(src map[string]interface{}) []byte {
	var lines []string

	for key, value := range src {
		var raw string

		switch t := value.(type) {
		case []int:
			raw = strings.Trim(strings.Replace(fmt.Sprint(value), " ", ",", -1), "[]")
		case []string:
			raw = strings.Join(value.([]string), ",")
		case map[string]interface{}:
			raw = ""
			fmt.Println("WARN: This type", t, "will be ignored")
		default:
			raw = fmt.Sprintf("%v", value)
		}
		lines = append(lines, fmt.Sprintf("export %s=%s", key, raw))
	}

	content := strings.Join(lines, "\n")
	return []byte(content)
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

		content[v[0]] = value
	}

	return content, nil
}
