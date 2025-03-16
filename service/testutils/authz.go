package testutils

import "user-service/authz"

type Datastore interface {
	Get(key string) interface{}
	Set(key string, val interface{}) error
}

type TestAuthZService struct {
	store Datastore
}

// Sample initialization
func InitTestAuthZService(db Datastore) *TestAuthZService {
	return &TestAuthZService{store: db}
}

func (s *TestAuthZService) IsAuthorized(role string, resource string, permission string, conditions interface{}) (bool, error) {
	return authz.IsRoleAuthorized(s.store, role, resource, permission, conditions)
}
