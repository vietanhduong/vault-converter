package util

import (
	"crypto/sha1"
	"encoding/hex"
)

func SHA1(src []byte) string {
	h := sha1.New()
	h.Write(src)
	return hex.EncodeToString(h.Sum(nil))
}