package vault

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/vietanhduong/vault-converter/pkg/cerror"
	"github.com/vietanhduong/vault-converter/pkg/util/util"
	"net/http"
	"strings"
)

type Vault struct {
	Address string
	Token   string
}

func New(vaultAddr, clientToken string) *Vault {
	return &Vault{Address: vaultAddr, Token: clientToken}
}

// Read read specified secret path and return a map
func (v *Vault) Read(secretPath string) (map[string]interface{}, error) {
	secretURL := util.JoinURL(fmt.Sprintf("%s/v1", v.Address), secretPath)

	req, err := http.NewRequest(http.MethodGet, secretURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Vault: Init request to read secret failed")
	}

	req.Header.Set("X-Vault-Token", v.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
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
