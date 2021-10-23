package pull

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)
// Converter interface. (PULL)
// Support converting from JSON (Vault) to specified format.
type Converter interface {
	Convert(src map[string]interface{}, output string) error
}

// NewConverter get a converter base on input format
func NewConverter(format string) (Converter, error) {
	switch strings.ToLower(format) {
	case "tfvars":
		return NewTfvars(), nil
	case "env":
		return NewEnv(), nil
	default:
		return nil, errors.New(fmt.Sprintf("`%s` is not yet supported.", format))
	}
}
