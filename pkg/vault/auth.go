package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/vietanhduong/vault-converter/pkg/cerror"
)

type Auth struct {
	VaultAddr string
	Username  string
	password  string
	client    HttpClient
}

func NewAuth(vaultAddr, username, password string) *Auth {
	return &Auth{
		VaultAddr: vaultAddr,
		Username:  username,
		password:  password,
		client:    &http.Client{},
	}
}

// Login authenticate the input user to vault server.
// In the current version, the auth pass of 'userpass'
// has been fixed with 'userpass/'
func (a *Auth) Login() (string, error) {
	loginURL := fmt.Sprintf("%s/v1/auth/userpass/login/%s", a.VaultAddr, a.Username)

	payload, _ := json.Marshal(&AuthPayload{Password: a.password})

	req, err := http.NewRequest(http.MethodPost, loginURL, bytes.NewBuffer(payload))
	if err != nil {
		return "", errors.Wrap(err, "Auth: Init request to login failed")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "Auth: Request login failed")
	}

	var ret *Response
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return "", errors.Wrap(err, "Auth: Read response body failed")
	}

	if resp.StatusCode != http.StatusOK {
		msg := cerror.DefaultErrorMsg(resp.StatusCode)
		if len(ret.Errors) > 0 {
			msg = ret.Errors[0]
		}
		return "", errors.New(fmt.Sprintf("[%d] Auth: %s", resp.StatusCode, strings.Title(msg)))
	}

	return ret.Auth.ClientToken, nil
}
