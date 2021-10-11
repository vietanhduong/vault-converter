package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

type ClientMock struct {
	DoFunc func(*http.Request) (*http.Response, error)
}

func (c *ClientMock) Do(r *http.Request) (*http.Response, error) {
	return c.DoFunc(r)
}

func TestAuth_Login(t *testing.T) {
	
	var (
		token     = "token"
		vaultAddr = "http://127.0.0.1:8200"
		username  = "root"
		password  = "123456"
	)

	t.Run("With success case: return token", func(tc *testing.T) {
		mr, _ := json.Marshal(&Response{Auth: &ResponseAuth{ClientToken: token}})
		client := &ClientMock{
			DoFunc: func(request *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader(mr)),
				}, nil
			},
		}
		a := NewAuth(vaultAddr, username, password)
		a.client = client
		tokenActual, err := a.Login()
		assert.NoError(tc, err)
		assert.Equal(tc, token, tokenActual)
	})

	t.Run("With success case: return token", func(tc *testing.T) {
		var errMsg = "Test error"
		client := &ClientMock{
			DoFunc: func(request *http.Request) (*http.Response, error) {
				return nil, errors.New(errMsg)
			},
		}
		a := NewAuth(vaultAddr, username, password)
		a.client = client
		_, err := a.Login()
		assert.EqualError(tc, err, fmt.Sprintf("Auth: Request login failed: %s", errMsg))
	})

	t.Run("With error case: failure with read response", func(tc *testing.T) {
		client := &ClientMock{
			DoFunc: func(request *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(`ok`))),
				}, nil
			},
		}
		a := NewAuth(vaultAddr, username, password)
		a.client = client
		_, err := a.Login()
		assert.EqualError(tc, err, "Auth: Read response body failed: invalid character 'o' looking for beginning of value")
	})

	t.Run("With error case: failure with response code ne 200", func(tc *testing.T) {
		mr, _ := json.Marshal(&Response{Errors: []string{"invalid username or password"}})
		client := &ClientMock{
			DoFunc: func(request *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 401,
					Body:       ioutil.NopCloser(bytes.NewReader(mr)),
				}, nil
			},
		}
		a := NewAuth(vaultAddr, username, password)
		a.client = client
		_, err := a.Login()
		assert.EqualError(tc, err, "[401] Auth: Invalid Username Or Password")
	})

}
