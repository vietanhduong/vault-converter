package converter

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type Converter interface {
	Convert(src map[string]interface{}, output string) error
}

func NewConverter(format string) (Converter, error) {
	switch strings.ToLower(format) {
	case "tfvars":
		return NewTfvars(), nil
	default:
		return nil, errors.New(fmt.Sprintf("`%s` is not yet supported.", format))
	}
}
