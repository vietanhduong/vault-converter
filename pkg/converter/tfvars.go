package converter

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/hcl/hcl/printer"
	jsonParser "github.com/hashicorp/hcl/json/parser"
	"github.com/pkg/errors"
	osext "github.com/vietanhduong/vault-converter/pkg/util/os"
	"os"
)

type tfvars struct {
}

func NewTfvars() Converter {
	return &tfvars{}
}

func (t *tfvars) Convert(src map[string]interface{}, output string) error {
	input, err := json.Marshal(src)
	if err != nil {
		return errors.Wrap(err, "Convert: Convert map to json failed")
	}

	ast, err := jsonParser.Parse(input)
	if err != nil {
		return errors.Wrap(err, "Convert: Unable to parse json")
	}

	if err = osext.MkdirP(output); err != nil {
		return errors.Wrap(err, "Convert: Create folder failed")
	}

	f, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Convert: Open %s failed", output))
	}
	defer f.Close()

	return printer.Fprint(f, ast)
}
