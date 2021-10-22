package push

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

// Converter interface. (PUSH)
// Support convert to JSON from specified format.
type Converter interface {
	Convert(src []byte) (map[string]interface{}, error)
}

// NewConverter get a converter by input format
func NewConverter(format string) (Converter, error) {
	switch strings.ToLower(format) {
	case ".tfvars":
		return NewTfvars(), nil
	case ".env":
		return NewEnv(), nil
	default:
		return nil, errors.New(fmt.Sprintf("`%s` is not yet supported.", format))
	}
}
