package os

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/vietanhduong/vault-converter/pkg/util/output"
	"io/ioutil"
	"os"
	"path"
)

func Cat(path string) (string, error) {
	file, err := os.Stat(path)
	if err != nil {
		return "", errors.Wrap(err, "Cat file error")
	}
	if file.IsDir() {
		return "", errors.New(fmt.Sprintf("%s: Path is a directory", path))
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.Wrap(err, "Cat: Read file error")
	}

	return string(b), nil
}

func Write(content []byte, output string) error {
	// Make sure output path exist
	dir := path.Dir(output)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Create directory with Mode 0755
		if err = os.MkdirAll(dir, 0755); err != nil {
			return errors.Wrap(err, "Write: Mkdir failed")
		}
	}

	// Write to file (No append)
	f, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "Write: OpenFile failed")
	}
	defer f.Close()

	if _, err = f.Write(content); err != nil {
		return errors.Wrap(err, "Write: Write content failed")
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
