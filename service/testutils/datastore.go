package testutils

import (
	"user-service/authz"
	"user-service/users"
)

type TestStore struct {
	Data map[string]interface{}
}

func (s *TestStore) Get(key string) interface{} {
	return s.Data[key]
}

func (s *TestStore) Set(key string, val interface{}) error {
	s.Data[key] = val
	return nil
}

// Sample initialization
func InitTestStore() *TestStore {
	sampleUsers := users.UsersInDB{
		"client_user": {
			ID:   "client_user",
			Name: "john.doe",
		},
		"user1": {
			ID:   "user1",
			Name: "john.doe",
		},
		"user2": {
			ID:   "user2",
			Name: "jane.smith",
		},
	}

	userRoles := users.UserRoles{
		users.UserID("client_user"): []authz.Role{authz.RoleViewer},
	}

	// sample rbac state, with just one role viewer with permission to read all users
	rbac := authz.RbacInDB{
		authz.RoleViewer: authz.AccessRights{
			Role:        authz.RoleViewer,
			Resource:    authz.ResourceUser,
			Permissions: []authz.Permission{authz.PermissionRead},
			Conditions: authz.Conditions{
				authz.CondKeyResourceID: "*",
			},
		},
	}

	return &TestStore{
		Data: map[string]interface{}{
			"users":      sampleUsers,
			"user_roles": userRoles,
			"rbac":       rbac,
		},
	}
}
