package hcl

import (
	"encoding/json"
	"fmt"
	hclv2 "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/vietanhduong/vault-converter/pkg/util/util"
	"github.com/zclconf/go-cty/cty"
)

type Hcl struct {
	parser *hclparse.Parser
}

func New() *Hcl {
	return &Hcl{
		parser: hclparse.NewParser(),
	}
}

// ToJSON convert input source to map.
// src should have HCL format.
func (h *Hcl) ToJSON(src []byte) (map[string]interface{}, error) {

	content, diags := h.parser.ParseHCL(src, util.SHA1(src))
	if diags.HasErrors() {
		return nil, diags
	}

	vars, diags := parseVarsBody(content.Body)
	if diags.HasErrors() {
		return nil, diags
	}

	values := make(map[string]interface{})

	for key, variable := range vars {
		values[key] = parseCtyValue(variable)
	}
	return values, nil
}

// ToHCL convert input source to HCL content
func (h *Hcl) ToHCL(src map[string]interface{}) ([]byte, error) {
	input, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}

	f, _ := h.parser.ParseJSON(input, util.SHA1(input))

	vars, diags := parseVarsBody(f.Body)
	if diags.HasErrors() {
		return nil, diags
	}

	// Create an empty file
	outfile := hclwrite.NewEmptyFile()
	body := outfile.Body()

	// Write content to body
	for name, value := range vars {
		body.SetAttributeValue(name, value)
		body.AppendNewline()
	}

	// Format before return content
	return hclwrite.Format(outfile.Bytes()), nil
}

// parseVarsBody read the file hcl file content
// and get all attributes. Since we just need to process
// the `tfvars` file, so just get the attributes is enough
func parseVarsBody(body hclv2.Body) (map[string]cty.Value, hclv2.Diagnostics) {
	attrs, diags := body.JustAttributes()
	if attrs == nil || diags.HasErrors() {
		return nil, diags
	}

	values := make(map[string]cty.Value, len(attrs))
	for name, attr := range attrs {
		val, valDiags := attr.Expr.Value(nil)
		diags = append(diags, valDiags...)
		values[name] = val
	}
	return values, diags
}

// TODO: Move to global function
func parseCtyValue(val cty.Value) interface{} {
	switch {
	case !val.IsKnown():
		panic("cannot produce tokens for unknown value")

	case val.IsNull():
		return nil

	case val.Type() == cty.Bool:
		return val.True()

	case val.Type() == cty.Number:
		return val.AsBigFloat()

	case val.Type() == cty.String:
		return val.AsString()

	case val.Type().IsListType() || val.Type().IsSetType() || val.Type().IsTupleType():
		var ls []interface{}
		for it := val.ElementIterator(); it.Next(); {
			_, eVal := it.Element()
			ls = append(ls, parseCtyValue(eVal))
		}
		return ls

	case val.Type().IsMapType() || val.Type().IsObjectType():
		m := make(map[string]interface{})

		for it := val.ElementIterator(); it.Next(); {
			eKey, eVal := it.Element()
			m[eKey.AsString()] = parseCtyValue(eVal)
		}

		return m
	default:
		panic(fmt.Sprintf("cannot procude value for %#v", val))
	}
}
