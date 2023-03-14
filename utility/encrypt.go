package utility

import (
	"crypto/md5"
	"encoding/hex"
)

// Md5Hash 算法
func Md5Hash(data string) string {
	Md5 := md5.New()
	Md5.Write([]byte(data))
	dataBytes := Md5.Sum(nil)
	return hex.EncodeToString(dataBytes)
}
