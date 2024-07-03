package helper

import (
	"crypto/md5"
	"encoding/hex"
)

func CreateHash(key string) []byte {

	hasher := md5.New()
	hasher.Write([]byte(key))
	return []byte(hex.EncodeToString(hasher.Sum(nil)))

}
