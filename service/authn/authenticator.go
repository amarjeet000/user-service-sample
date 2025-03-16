// Package authn handles the token generation and validation for the authentication purposes.
package authn

import (
	"errors"
	"sync"
	"user-service/config"
)

// Hardcoded clientId, to be used as audience during token generation.
// A clientId is generally issued by the auth server when a client registers with the server.
const clientId = "client1"

// Service implements Authenticator interface
type Service struct {
	sync.RWMutex
	Name   string
	Secret string // To be used when using HMAC signing method for token
	Cfg    *config.Config
}

// InitService initializes the Service with hardcoded values of name and secret.
// The secret is included here to enable testing for HMAC signed token mechanism.
func InitService() *Service {
	return &Service{
		Name:   "platform/user-service",
		Secret: "mysupersecret",
	}
}

// ClientClaims is just a sample claims data structure.
// It does not offer much value in this sample user-service.
type ClientClaims struct {
	// The UserID may indicate a human-user or a machine-user.
	UserID string `json:"user_id"`
}

func (s *Service) GenerateToken(id string) (string, error) {
	// We use mutext to remove any rare possibility of two tokens having the same properties because of concurrent calls.
	// This lock only applies to writes.
	// The goal is to not block any readers of the Service object.
	s.Lock()
	defer s.Unlock()
	switch s.Cfg.SigningMethod {
	case "rsa":
		return GenerateRSASignedToken(s.Cfg, id, s.Name)
	case "hmac":
		return GenerateHMACSignedToken(id, s.Name, s.Secret)
	default:
		return "", errors.New("invalid signing-method")
	}
}

func (s *Service) ValidateToken(token string) (string, error) {
	switch s.Cfg.SigningMethod {
	case "rsa":
		return ValidateRSASignedToken(s.Cfg, token, s.Name)
	case "hmac":
		return ValidateHMACSignedToken(token, s.Name, s.Secret)
	default:
		return "", errors.New("invalid signing-method")
	}
}
