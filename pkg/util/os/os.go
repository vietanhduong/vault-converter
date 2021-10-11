package os

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/vietanhduong/vault-converter/pkg/util/output"
	"github.com/vietanhduong/vault-converter/pkg/util/util"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func Cat(path string) ([]byte, error) {
	file, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrap(err, "Cat file error")
	}
	if file.IsDir() {
		return nil, errors.New(fmt.Sprintf("%s: Path is a directory", path))
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "Cat: Read file error")
	}

	return b, nil
}

func Write(content []byte, output string) error {
	// Make sure output path exist
	if _, err := os.Stat(output); os.IsNotExist(err) {
		dir := path.Dir(output)
		// Create directory with Mode 0755
		if err = os.MkdirAll(dir, 0755); err != nil {
			return errors.Wrap(err, "Write: Mkdir failed")
		}
	}

	// Write to file (No append)
	f, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		return errors.Wrap(err, "Write: OpenFile failed")
	}
	defer f.Close()

	if _, err = f.Write(content); err != nil {
		return errors.Wrap(err, "Write: Write content failed")
	}

	return nil
}

func MkdirP(input string) error {
	if _, err := os.Stat(input); os.IsNotExist(err) {
		dir := path.Dir(input)
		// Create directory with Mode 0755
		if err = os.MkdirAll(dir, 0755); err != nil {
			return errors.Wrap(err, "Write: Mkdir failed")
		}
	}
	return nil
}

func HomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		output.Printf("WARN: $HOME is not defined")
	}
	return home
}

func GetExtension(path string) string {
	ext := filepath.Ext(path)
	if util.IsNullOrEmpty(ext) {
		return ""
	}
	return ext
}
