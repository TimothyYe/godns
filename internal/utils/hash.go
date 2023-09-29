package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// GetMD5Hash returns the MD5 hash of the input string.
func GetMD5Hash(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}
