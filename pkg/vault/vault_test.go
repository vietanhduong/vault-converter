package vault

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVault_Read(t *testing.T) {
	var (
		token     = "token"
		vaultAddr = "http://127.0.0.1:8200"
	)

	t.Run("With success case: return map contain values", func(tc *testing.T) {
		var secretPath = "test/data/cluster"
		values := map[string]interface{}{
			"pass": true,
		}
		mr, _ := json.Marshal(&Response{Data: &ResponseData{Data: values}})

		v := &vault{
			addr:  vaultAddr,
			token: token,
		}
		v.client = &ClientMock{
			DoFunc: func(request *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader(mr)),
				}, nil
			},
		}

		data, err := v.Read(secretPath)
		assert.NoError(tc, err)
		assert.Equal(tc, len(data), len(values))
		assert.Equal(tc, data["pass"], values["pass"])
	})

	t.Run("With error case: failure with do request", func(tc *testing.T) {
		var secretPath = "test/data/cluster"

		v := New("", token)
		_, err := v.Read(secretPath)
		assert.Error(tc, err)
		assert.Contains(tc, err.Error(), "vault: Request to read secret failed")
	})

	t.Run("With error case: failure with marshal response", func(tc *testing.T) {
		var secretPath = "test/data/cluster"
		v := &vault{
			addr:  vaultAddr,
			token: token,
		}

		v.client = &ClientMock{
			DoFunc: func(request *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(`ok`))),
				}, nil
			},
		}

		_, err := v.Read(secretPath)
		assert.Error(tc, err)
		assert.Contains(tc, err.Error(), "vault: Read response body failed")
	})

	t.Run("With error case: failure with response code ne 200", func(tc *testing.T) {
		var secretPath = "test/data/cluster"
		mr, _ := json.Marshal(&Response{Errors: []string{"* Permission denied"}})
		v := &vault{
			addr:  vaultAddr,
			token: token,
		}
		v.client = &ClientMock{
			DoFunc: func(request *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 403,
					Body:       ioutil.NopCloser(bytes.NewReader(mr)),
				}, nil
			},
		}
		_, err := v.Read(secretPath)
		assert.Error(tc, err)
		assert.Contains(tc, err.Error(), "[403] vault: * Permission Denied")
	})
}

func TestVault_Write(t *testing.T) {
	var (
		token     = "token"
		vaultAddr = "http://127.0.0.1:8200"
	)

	t.Run("With success case: return no error", func(tc *testing.T) {
		var secretPath = "test/data/cluster"
		mr, _ := json.Marshal(&Response{Data: &ResponseData{}})

		v := &vault{
			addr:  vaultAddr,
			token: token,
		}
		v.client = &ClientMock{
			DoFunc: func(request *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader(mr)),
				}, nil
			},
		}

		err := v.Write(secretPath, map[string]interface{}{})
		assert.NoError(tc, err)
	})

	t.Run("With error case: failure with do request", func(tc *testing.T) {
		var secretPath = "test/data/cluster"

		v := New("", token)
		err := v.Write(secretPath, map[string]interface{}{})
		assert.Error(tc, err)
		assert.Contains(tc, err.Error(), "vault: Request write secret failed")
	})

	t.Run("With error case: failure with marshal response", func(tc *testing.T) {
		var secretPath = "test/data/cluster"
		v := &vault{
			addr:  vaultAddr,
			token: token,
		}
		v.client = &ClientMock{
			DoFunc: func(request *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(`ok`))),
				}, nil
			},
		}

		err := v.Write(secretPath, map[string]interface{}{})
		assert.Error(tc, err)
		assert.Contains(tc, err.Error(), "vault: Read response body failed")
	})

	t.Run("With error case: failure with response code ne 200", func(tc *testing.T) {
		var secretPath = "test/data/cluster"
		mr, _ := json.Marshal(&Response{Errors: []string{"* Permission denied"}})
		v := &vault{
			addr:  vaultAddr,
			token: token,
		}
		v.client = &ClientMock{
			DoFunc: func(request *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 403,
					Body:       ioutil.NopCloser(bytes.NewReader(mr)),
				}, nil
			},
		}
		err := v.Write(secretPath, map[string]interface{}{})
		assert.Error(tc, err)
		assert.Contains(tc, err.Error(), "[403] vault: * Permission Denied")
	})

}

func TestVault_List(t *testing.T) {
	var (
		token     = "token"
		vaultAddr = "http://127.0.0.1:8200"
	)

	t.Run("With success case: path is empty", func(tc *testing.T) {
		var mr *Response
		v := &vault{
			addr:  vaultAddr,
			token: token,
		}
		v.client = &ClientMock{
			DoFunc: func(request *http.Request) (*http.Response, error) {
				switch request.URL.Path {
				case "/v1/sys/internal/ui/mounts":
					mr = &Response{
						Data: &ResponseData{
							Secret: map[string]interface{}{
								"data/": nil,
								"sys/":  nil,
							},
						},
					}
				case "/v1/data/metadata":
					mr = &Response{
						Data: &ResponseData{
							Keys: []string{"secrets", "aws/"},
						},
					}
				case "/v1/data/metadata/aws":
					mr = &Response{
						Data: &ResponseData{
							Keys: []string{"content"},
						},
					}
				}

				b, _ := json.Marshal(mr)
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader(b)),
				}, nil
			},
		}

		secrets, err := v.List("", true)
		assert.NoError(tc, err)
		assert.Equal(tc, len(secrets), 2)
		assert.Contains(tc, secrets, "data/secrets")
		assert.Contains(tc, secrets, "data/aws/content")
	})

	t.Run("With success case: path is not empty", func(tc *testing.T) {
		var mr *Response
		v := &vault{
			addr:  vaultAddr,
			token: token,
		}
		v.client = &ClientMock{
			DoFunc: func(request *http.Request) (*http.Response, error) {
				switch request.URL.Path {
				case "/v1/sys/internal/ui/mounts":
					mr = &Response{
						Data: &ResponseData{
							Secret: map[string]interface{}{
								"data/": nil,
								"test/": nil,
								"sys/":  nil,
							},
						},
					}
				case "/v1/data/metadata":
					mr = &Response{
						Data: &ResponseData{
							Keys: []string{"secrets", "aws/"},
						},
					}
				case "/v1/data/metadata/aws":
					mr = &Response{
						Data: &ResponseData{
							Keys: []string{"content"},
						},
					}
				}

				b, _ := json.Marshal(mr)
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader(b)),
				}, nil
			},
		}

		secrets, err := v.List("data/aws", true)
		assert.NoError(tc, err)
		assert.Equal(tc, len(secrets), 1)
		assert.Contains(tc, secrets, "data/aws/content")
	})

	t.Run("With error case: invalid path", func(tc *testing.T) {
		var mr *Response
		v := &vault{
			addr:  vaultAddr,
			token: token,
		}
		v.client = &ClientMock{
			DoFunc: func(request *http.Request) (*http.Response, error) {
				switch request.URL.Path {
				case "/v1/sys/internal/ui/mounts":
					mr = &Response{
						Data: &ResponseData{
							Secret: map[string]interface{}{
								"data/": nil,
								"sys/":  nil,
							},
						},
					}
				}

				b, _ := json.Marshal(mr)
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader(b)),
				}, nil
			},
		}

		_, err := v.List("test", true)
		assert.Error(tc, err)
		assert.Contains(tc, err.Error(), "Path 'test' no longer exist")
	})

	t.Run("With error case: fetch root failed", func(tc *testing.T) {
		var mr *Response
		v := &vault{
			addr:  vaultAddr,
			token: token,
		}
		v.client = &ClientMock{
			DoFunc: func(request *http.Request) (*http.Response, error) {
				switch request.URL.Path {
				case "/v1/sys/internal/ui/mounts":
					mr = &Response{
						Errors: []string{"error"},
					}
				}

				b, _ := json.Marshal(mr)
				return &http.Response{
					StatusCode: 400,
					Body:       ioutil.NopCloser(bytes.NewReader(b)),
				}, nil
			},
		}

		_, err := v.List("data", true)
		assert.Error(tc, err)
		assert.Contains(tc, err.Error(), "List secrets failed")
	})

}
