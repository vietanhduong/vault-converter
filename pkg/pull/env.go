package pull

import (
	"fmt"
	"github.com/pkg/errors"
	envlib "github.com/vietanhduong/vault-converter/pkg/env"
	osext "github.com/vietanhduong/vault-converter/pkg/util/os"
)

type env struct {
	envLib *envlib.Env
}

func NewEnv() Converter {
	return &env{
		envLib: envlib.NewEnv(),
	}
}

// Convert a source to ENV file
// src input is a map. It should be a JSON format
// output should be an absolute path
func (e *env) Convert(src map[string]interface{}, output string) error {

	content := e.envLib.ToENV(src)

	if err := osext.MkdirP(output); err != nil {
		return errors.Wrap(err, "Pull: Create folder failed")
	}

	if err := osext.Write(content, output); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Pull: Write to %s failed", output))
	}

	return nil
}
