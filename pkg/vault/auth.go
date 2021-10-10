package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/vietanhduong/vault-converter/pkg/util/os"
	"net/http"
	"strings"
)

type Auth struct {
	VaultAddr string
	Username  string
	password  string
}

func NewAuth(vaultAddr, username, password string) *Auth {
	return &Auth{
		VaultAddr: vaultAddr,
		Username:  username,
		password:  password,
	}
}

func (a *Auth) Login() error {
	loginURL := fmt.Sprintf("%s/v1/auth/userpass/login/%s", a.VaultAddr, a.Username)
	payload, err := json.Marshal(&AuthPayload{Password: a.password})
	if err != nil {
		return errors.Wrap(err, "Auth: Marshal auth payload failed")
	}
	req, err := http.NewRequest(http.MethodPost, loginURL, bytes.NewBuffer(payload))
	if err != nil {
		return errors.Wrap(err, "Auth: Init request to login failed")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "Auth: Request login failed")
	}

	var ret *Response
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return errors.Wrap(err, "Auth: Read response body failed")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("[%d] Auth: %s", resp.StatusCode, strings.Title(ret.Errors[0])))
	}

	if err = os.Write([]byte(ret.Auth.ClientToken), DefaultTokenPath); err != nil {
		return errors.Wrap(err, "Login: Write token to file failed")
	}
	return nil
}
