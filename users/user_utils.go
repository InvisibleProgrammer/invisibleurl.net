package users

import (
	"crypto/rand"
	"encoding/base64"
)

func generateRandomState() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
