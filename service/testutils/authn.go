package testutils

import "user-service/authn"

type TestAuthNService struct {
	Name   string
	Secret string
}

func InitTestAuthNService() *TestAuthNService {
	return &TestAuthNService{
		Name:   "user-service",
		Secret: "mysupersecret",
	}
}

func (s *TestAuthNService) GenerateToken(id string) (string, error) {
	return authn.GenerateHMACSignedToken(id, s.Name, s.Secret)
}

func (s *TestAuthNService) ValidateToken(token string) (string, error) {
	return authn.ValidateHMACSignedToken(token, s.Name, s.Secret)
}
