package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// md5 hash

func Md5(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}
