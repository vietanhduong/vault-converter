package os

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestCat(t *testing.T) {
	// Setup test
	var content = "test_content"
	tmpDir := os.TempDir()
	f, _ := ioutil.TempFile(tmpDir, "test")
	defer os.Remove(f.Name())
	_, _ = f.Write([]byte(content))

	t.Run("With success case: return content file", func(tc *testing.T) {
		c, err := Cat(f.Name())
		assert.NoError(tc, err)
		assert.Equal(tc, string(c), content)
	})

	t.Run("With error case: failure with file not exist", func(tc *testing.T) {
		_, err := Cat(f.Name() + "test")
		assert.Error(tc, err)
		assert.Contains(tc, err.Error(), "Cat file error")
	})

	t.Run("With error case: failure with input path is directory", func(tc *testing.T) {
		_, err := Cat(tmpDir)
		assert.Error(tc, err)
		assert.Contains(tc, err.Error(), "Path is a directory")
	})

	t.Run("With error case: failure with read content file", func(tc *testing.T) {
		_ = os.Chmod(f.Name(), 0100)
		_, err := Cat(f.Name())
		assert.Error(tc, err)
		assert.Contains(tc, err.Error(), "Cat file error")
	})
}

func TestWrite(t *testing.T) {

}