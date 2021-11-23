package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/vietanhduong/vault-converter/pkg/cerror"
	"github.com/vietanhduong/vault-converter/pkg/util/os"
	"github.com/vietanhduong/vault-converter/pkg/util/util"
)

var DefaultTokenPath = os.HomeDir() + "/.vault_converter/token"

type HttpClient interface {
	Do(r *http.Request) (*http.Response, error)
}

type Vault struct {
	Address string
	Token   string
	client  HttpClient
}

func New(vaultAddr, clientToken string) *Vault {
	return &Vault{
		Address: vaultAddr,
		Token:   clientToken,
		client:  &http.Client{},
	}
}

// Read read specified secret path and return a map
func (v *Vault) Read(secretPath string) (map[string]interface{}, error) {
	secretURL := util.JoinURL(fmt.Sprintf("%s/v1", v.Address), secretPath)

	req, err := http.NewRequest(http.MethodGet, secretURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Vault: Init request to read secret failed")
	}

	req.Header.Set("X-Vault-Token", v.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := v.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Vault: Request to read secret failed")
	}

	var ret *Response
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, errors.Wrap(err, "Vault: Read response body failed")
	}

	if resp.StatusCode != http.StatusOK {
		msg := cerror.DefaultErrorMsg(resp.StatusCode)
		if len(ret.Errors) > 0 {
			msg = ret.Errors[0]
		}
		return nil, errors.New(fmt.Sprintf("[%d] Vault: %s", resp.StatusCode, strings.Title(msg)))
	}

	return ret.Data.Data, nil
}

func (v *Vault) Write(secretPath string, values map[string]interface{}) error {
	secretURL := util.JoinURL(fmt.Sprintf("%s/v1", v.Address), secretPath)
	payloadObject := &SecretPayload{Data: values}

	payload, _:= json.Marshal(payloadObject)

	req, err := http.NewRequest(http.MethodPost, secretURL, bytes.NewBuffer(payload))
	if err != nil {
		return errors.Wrap(err, "Vault: Init request to create secret failed")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Vault-Token", v.Token)

	resp, err := v.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "Vault: Request write secret failed")
	}

	var ret *Response
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return errors.Wrap(err, "Vault: Read response body failed")
	}

	if resp.StatusCode != 200 {
		msg := cerror.DefaultErrorMsg(resp.StatusCode)
		if len(ret.Errors) > 0 {
			msg = ret.Errors[0]
		}
		return errors.New(fmt.Sprintf("[%d] Vault: %s", resp.StatusCode, strings.Title(msg)))
	}

	return nil
}
