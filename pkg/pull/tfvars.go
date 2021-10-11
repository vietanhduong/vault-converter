package pull

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/hcl/hcl/printer"
	jsonParser "github.com/hashicorp/hcl/json/parser"
	"github.com/pkg/errors"
	osext "github.com/vietanhduong/vault-converter/pkg/util/os"
	out "github.com/vietanhduong/vault-converter/pkg/util/output"
	"path/filepath"
	"regexp"
)

type tfvars struct {
}

func NewTfvars() Converter {
	return &tfvars{}
}

// Convert convert a source to a HCL file
// src is a map. It should be a JSON format
func (t *tfvars) Convert(src map[string]interface{}, output string) error {
	input, err := json.Marshal(src)
	if err != nil {
		return errors.Wrap(err, "Pull: Convert map to json failed")
	}

	ast, err := jsonParser.Parse(input)
	if err != nil {
		return errors.Wrap(err, "Pull: Unable to parse json")
	}

	if err = osext.MkdirP(output); err != nil {
		return errors.Wrap(err, "Pull: Create folder failed")
	}

	var buff []byte
	buffer := bytes.NewBuffer(buff)

	if err = printer.Fprint(buffer, ast); err != nil {
		return errors.Wrap(err, "Pull: Write buffer failed")
	}

	// There a issue in here
	// Ref: https://github.com/hashicorp/hcl/issues/233
	// When parsing JSON with github.com/hashicorp/hcl/json/parser
	// and printing HCL with github.com/hashicorp/hcl/hcl/printer,
	// quotes in function args are broken.
	// The temporary solution is to use regex to remove the double quotes.
	// But this will reduce performance.
	content := buffer.String()
	match := regexp.MustCompile(`"(\w+)"\s+=`)
	content = match.ReplaceAllString(content, "$1 =")

	// Format the content before write to file
	formattedContent, err := printer.Format([]byte(content))
	if err = osext.Write(formattedContent, output); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Pull: Write to %s failed", output))
	}

	absPath, _ := filepath.Abs(output)
	out.Printf("Generated output at %s", absPath)
	return nil
}
