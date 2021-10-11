package push

import (
	"github.com/hashicorp/hcl"
	"github.com/pkg/errors"
)

type tfvars struct {
}

func NewTfvars() Converter {
	return &tfvars{}
}

func (t *tfvars) Convert(src []byte) (map[string]interface{}, error) {
	ret := make(map[string]interface{})
	err := hcl.Unmarshal(src, &ret)
	if err != nil {
		return nil, errors.Wrap(err, "Push: Unmarshal source failed")
	}
	return ret, nil
}
