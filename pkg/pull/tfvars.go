package pull

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/vietanhduong/vault-converter/pkg/hcl"
	osext "github.com/vietanhduong/vault-converter/pkg/util/os"
)

type tfvars struct {
	hclConverter *hcl.Hcl
}

func NewTfvars() Converter {
	return &tfvars{
		hclConverter: hcl.New(),
	}
}

// Convert a source to HCL file
// src is a map. It should be a JSON format
// output should be an absolute path
func (t *tfvars) Convert(src map[string]interface{}, output string) error {

	content, err := t.hclConverter.ToHCL(src)
	if err != nil {
		return errors.Wrap(err, "Pull: Convert to HCL failed")
	}

	if err = osext.MkdirP(output); err != nil {
		return errors.Wrap(err, "Pull: Create folder failed")
	}

	if err = osext.Write(content, output); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Pull: Write to %s failed", output))
	}

	return nil
}
