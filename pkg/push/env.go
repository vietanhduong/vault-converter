package push

import (
	"github.com/pkg/errors"
	envlib "github.com/vietanhduong/vault-converter/pkg/env"
	"strings"
)

type env struct {
	envLib *envlib.Env
}

func NewEnv() Converter {
	return &env{
		envLib: envlib.NewEnv(),
	}
}

// Convert an ENV file to JSON
func (e *env) Convert(src []byte) (map[string]interface{}, error){

	lines := strings.Split(string(src), "\n")

	content, err := e.envLib.ToJSON(lines)
	if err != nil {
		return nil, errors.Wrap(err, "Push: Convert bytes to json failed")
	}

	return content, nil
}
