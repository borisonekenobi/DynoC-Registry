package commons

import (
	"crypto/sha256"
	"encoding/hex"
)

func CalculateSHA256(data []byte) (string, error) {
	hash := sha256.New()
	_, err := hash.Write(data)
	if err != nil {
		return "", err
	}

	sum := hash.Sum(nil)
	checksum := hex.EncodeToString(sum)

	return checksum, nil
}
