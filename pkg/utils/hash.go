package utils

import (
	"crypto/sha256"
)

func GetHash(v []byte) []byte {
	h := sha256.New()
	h.Write(v)
	return h.Sum(nil)
}
