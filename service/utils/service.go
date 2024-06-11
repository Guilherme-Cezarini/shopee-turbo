package utils

import (
	"crypto/rand"
)

type UtilsService struct {
}

func NewUtilsService() UtilsService {
	return UtilsService{}
}

// GenerateRandomKey generates a new random key
func (service UtilsService) GenerateRandomKey(length int) ([]byte, error) {
	key := make([]byte, length)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}
