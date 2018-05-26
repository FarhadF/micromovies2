package vault

import (
	"context"
	"golang.org/x/crypto/bcrypt"
)

//Business logic as interface
type Service interface {
	Hash(ctx context.Context, password string) (string, error)
	Validate(ctx context.Context, password string, hash string) (bool, error)
}

//implementation with empty struct (stateless)
type vaultService struct {
}

//constructor - we can later add initialization if needed
func NewService() Service {
	return vaultService{}
}

//implementation
func (vaultService) Hash(ctx context.Context, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

//implementation
func (vaultService) Validate(ctx context.Context, password string, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
