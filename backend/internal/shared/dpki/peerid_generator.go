package dpki

import (
	"crypto/sha1"
)

func GenerateDPKIPeerID(pubKey []byte) [20]byte {
	return sha1.Sum(pubKey)
}
