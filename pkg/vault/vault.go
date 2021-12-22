package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/vietanhduong/vault-converter/pkg/cerror"
	"github.com/vietanhduong/vault-converter/pkg/util/os"
	"github.com/vietanhduong/vault-converter/pkg/util/util"
)

var DefaultTokenPath = os.HomeDir() + "/.vault_converter/token"
var SkipWithPaths = map[string]bool{
	"cubbyhole": true,
	"identity":  true,
	"sys":       true,
}

type HttpClient interface {
	Do(r *http.Request) (*http.Response, error)
}

type Client interface {
	Read(secretPath string) (map[string]interface{}, error)
	Write(secretPath string, values map[string]interface{}) error
	List(secretPath string, recursive bool) ([]string, error)
}

type vault struct {
	addr   string
	token  string
	client HttpClient
}

func New(vaultAddr, clientToken string) Client {
	return &vault{
		addr:   vaultAddr,
		token:  clientToken,
		client: &http.Client{},
	}
}

// Read specified secret path and return a map
func (v *vault) Read(path string) (map[string]interface{}, error) {
	secretURL := util.JoinURL(fmt.Sprintf("%s/v1", v.addr), path)

	resp, err := v.makeRequest(http.MethodGet, secretURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "vault: Request to read secret failed")
	}

	var ret *Response
	if err = json.Unmarshal(resp.Body, &ret); err != nil {
		return nil, errors.Wrap(err, "vault: Read response body failed")
	}

	if resp.StatusCode == http.StatusOK {
		return ret.Data.Data, nil
	}

	return nil, handleError(resp.StatusCode, ret)
}

func (v *vault) Write(path string, values map[string]interface{}) error {
	secretURL := util.JoinURL(fmt.Sprintf("%s/v1", v.addr), path)
	payloadObject := &SecretPayload{Data: values}

	payload, _ := json.Marshal(payloadObject)

	resp, err := v.makeRequest(http.MethodPost, secretURL, payload)
	if err != nil {
		return errors.Wrap(err, "vault: Request write secret failed")
	}

	var ret *Response
	if err = json.Unmarshal(resp.Body, &ret); err != nil {
		return errors.Wrap(err, "vault: Read response body failed")
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return handleError(resp.StatusCode, ret)
}

func (v *vault) List(path string, recursive bool) ([]string, error) {
	var secrets []string

	roots, err := v.FindRoots()
	if err != nil {
		return nil, errors.Wrap(err, "vault: List secrets failed")
	}

	if path == "" {
		for _, root := range roots {
			secrets = append(secrets, v.FindSecrets(root, "", recursive)...)
		}

		return secrets, nil
	}

	if !isValidPath(roots, path) {
		return nil, fmt.Errorf("vault: Path '%s' no longer exist", path)
	}

	folders := strings.Split(strings.ToLower(path), "/")
	root := folders[0]
	var p string
	if len(folders) > 1 {
		p = strings.Join(folders[1:], "/")
	}

	return v.FindSecrets(root, p, recursive), nil
}

func (v *vault) FindRoots() ([]string, error) {
	url := util.JoinURL(fmt.Sprintf("%s/v1", v.addr), "sys/internal/ui/mounts")

	r, err := v.makeRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var resp *Response
	if err = json.Unmarshal(r.Body, &resp); err != nil {
		return nil, err
	}

	if r.StatusCode != 200 {
		return nil, handleError(r.StatusCode, resp)
	}

	var mounts []string
	for k := range resp.Data.Secret {
		path := strings.TrimSuffix(k, "/")
		if _, ok := SkipWithPaths[path]; ok {
			continue
		}
		mounts = append(mounts, path)
	}

	return mounts, nil
}

func (v *vault) FindSecrets(root, path string, recursive bool) []string {
	url := util.JoinURL(fmt.Sprintf("%s/v1", v.addr), root, "metadata", path)
	r, err := v.makeRequest("LIST", url, nil)
	if err != nil {
		return nil
	}

	var resp *Response
	if err = json.Unmarshal(r.Body, &resp); err != nil {
		return nil
	}

	if r.StatusCode != http.StatusOK {
		return nil
	}

	var secrets []string
	for _, key := range resp.Data.Keys {
		secret := strings.TrimSuffix(key, "/")
		if !strings.HasSuffix(key, "/") {
			secrets = append(secrets, util.JoinURL(root, path, secret))
			continue
		}
		if recursive {
			secrets = append(secrets, v.FindSecrets(root, util.JoinURL(path, secret), recursive)...)
		}
	}

	return secrets
}

func (v *vault) makeRequest(method, url string, payload []byte) (*HttpResponse, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-vault-token", v.token)

	resp, err := v.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &HttpResponse{
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}

func handleError(statusCode int, resp *Response) error {
	msg := cerror.DefaultErrorMsg(statusCode)
	if len(resp.Errors) > 0 {
		msg = resp.Errors[0]
	}
	return errors.New(fmt.Sprintf("[%d] vault: %s", statusCode, strings.Title(msg)))
}

func isValidPath(roots []string, path string) bool {
	folders := strings.Split(strings.ToLower(path), "/")
	for _, root := range roots {
		if root == folders[0] {
			return true
		}
	}
	return false
}
