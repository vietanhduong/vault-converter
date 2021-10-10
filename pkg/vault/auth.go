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



func New(vaultAddr, username, password string) *Vault {
	return &Vault{
		VaultAddr: vaultAddr,
		Username:  username,
		password:  password,
	}
}

func (v *Vault) Login() error {
	loginURL := fmt.Sprintf("%s/v1/auth/userpass/login/%s", v.VaultAddr, v.Username)
	payload, err := json.Marshal(&AuthPayload{Password: v.password})
	if err != nil {
		return errors.Wrap(err, "Vault: Marshal auth payload failed")
	}

	resp, err := http.Post(loginURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return errors.Wrap(err, "Vault: Request login failed")
	}

	var ret *AuthResponse
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return errors.Wrap(err, "Vault: Read response body failed")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Vault: %s", strings.Title(ret.Errors[0])))
	}

	if err = os.Write([]byte(ret.Auth.ClientToken), DefaultTokenPath); err != nil {
		return errors.Wrap(err, "Login: Write token to file failed")
	}
	return nil
}
