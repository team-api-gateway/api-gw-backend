package hashid

import (
	"crypto/sha1"
	"encoding/hex"
)

func GenerateHashId(value string) string {
	h := sha1.New()
	h.Write([]byte(value))
	sha := h.Sum(nil)

	shaStr := hex.EncodeToString(sha)
	return shaStr
}
