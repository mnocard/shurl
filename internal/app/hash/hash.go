package hash

import (
	"crypto/sha1"
	"encoding/hex"
)

func GetHash(b []byte) string {
	h := sha1.New()
	h.Write(b)
	sha := hex.EncodeToString(h.Sum(nil))
	return sha[0:8]
}
