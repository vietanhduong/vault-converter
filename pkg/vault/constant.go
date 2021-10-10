package vault

import (
	"github.com/vietanhduong/vault-converter/pkg/util/os"
)

var DefaultTokenPath = os.HomeDir() + "/.vault_converter/token"