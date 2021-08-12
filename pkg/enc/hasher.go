package enc

import (
	"crypto/sha256"
	"fmt"
)

func Hashed(s string) string {
	hasher := sha256.New()
	hasher.Write([]byte(s))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
