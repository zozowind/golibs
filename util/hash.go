package util

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"io"
)

//Sha1 sha1 hash
func Sha1(str string) string {
	h := sha1.New()
	io.WriteString(h, str)
	return hex.EncodeToString(h.Sum(nil))
}

//Sha256 sha256 hash
func Sha256(str string) string {
	h := sha256.New()
	io.WriteString(h, str)
	return hex.EncodeToString(h.Sum(nil))
}

//Md5 md5 hash
func Md5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return hex.EncodeToString(h.Sum(nil))
}
