package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type HashingInterface interface {
	HashField(text string) string
}

type HashingConfig struct {
	First  string
	Second string
	Third  string
	Fourth string
}

type hashing struct {
	config HashingConfig
}

func NewHashing(config HashingConfig) (HashingInterface, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return &hashing{
		config: config,
	}, nil
}

func (h *hashing) HashField(text string) string {
	if strings.TrimSpace(text) == "" {
		return ""
	}

	methods := []string{h.config.First, h.config.Second, h.config.Fourth}
	response := text
	for _, method := range methods {
		h := hmac.New(sha256.New, []byte(method))
		h.Write([]byte(response))

		response = hex.EncodeToString(h.Sum(nil))
	}

	return response
}

func validateConfig(config HashingConfig) error {
	if config.First == "" {
		return errors.New("missing hashing key 1")
	}
	if config.Second == "" {
		return errors.New("missing hashing key 2")
	}
	if config.Third == "" {
		return errors.New("missing hashing key 3")
	}
	if config.Fourth == "" {
		return errors.New("missing hashing key 4")
	}

	return nil
}

func HashPasswordWithBcrypt(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHashWihBcrypt(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
