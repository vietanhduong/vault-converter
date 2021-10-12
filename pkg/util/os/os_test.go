package os

import (
	"fmt"
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
	var content = "test_content"
	tmpDir := os.TempDir()
	defer os.Remove(tmpDir)

	t.Run("With success case: write content without return error", func(tc *testing.T) {
		p := tmpDir + "/test_write1/test.txt"
		err := Write([]byte(content), p)
		assert.NoError(tc, err)
		c, err := Cat(p)
		assert.NoError(tc, err)
		assert.Equal(tc, content, string(c))
	})

	t.Run("With error case: failure with output is a dir", func(tc *testing.T) {
		err := Write([]byte(content), tmpDir)
		assert.Error(tc, err)
		assert.Contains(tc, err.Error(), fmt.Sprintf("Write: %s is a directory", tmpDir))
	})

	t.Run("With error case: failure with can not open file", func(tc *testing.T) {
		f, _ := ioutil.TempFile(tmpDir, "test")
		_ = os.Chmod(f.Name(), 0200)
		err := Write([]byte(content), f.Name())
		assert.Error(tc, err)
		assert.Contains(tc, err.Error(), "Write: OpenFile failed")
	})
}

func TestHomeDir(t *testing.T) {
	t.Run("With success case: return home dir", func(tc *testing.T) {
		h := "/usr"
		_ = os.Setenv("HOME", h)

		home := HomeDir()
		assert.Equal(tc, h, home)
	})

	t.Run("With failed case: return empty home path", func(tc *testing.T) {
		_ = os.Setenv("HOME", "")
		home := HomeDir()
		assert.Equal(tc, "", home)
	})
}

func TestGetExtension(t *testing.T) {
	t.Run("With success case: return extension", func(tc *testing.T) {
		p := "test.txt"
		expected := ".txt"
		a := GetExtension(p)
		assert.Equal(tc, expected, a)
	})

	t.Run("With failed case: return empty", func(tc *testing.T) {
		assert.Equal(tc, GetExtension(""), "")
	})
}