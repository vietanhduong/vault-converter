package push

import (
	"github.com/pkg/errors"
	"github.com/vietanhduong/vault-converter/pkg/hcl"
)

type tfvars struct {
	hclConverter *hcl.Hcl
}

func NewTfvars() Converter {
	return &tfvars{
		hclConverter: hcl.New(),
	}
}

func (t *tfvars) Convert(src []byte) (map[string]interface{}, error) {
	content, err := t.hclConverter.ToJSON(src)
	if err != nil {
		return nil, errors.Wrap(err, "Push: Parse content to JSON failed")
	}
	return content, nil
}
