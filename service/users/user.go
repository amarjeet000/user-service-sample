package users

import (
	"errors"
	"user-service/authz"
	"user-service/commons"
)

type UserID string

const (
	UserIdInReqCtx UserID = "id"
)

type User struct {
	ID   UserID `json:"id"`
	Name string `json:"username"`
}

type UserRoles map[UserID][]authz.Role // a map of user-id and array of roles

type UsersInDB map[UserID]User

func FetchUsersFilterOne(db commons.Datastore, userId string) ([]User, error) {
	usersInDB, ok := db.Get("users").(UsersInDB)
	if !ok {
		return nil, errors.New("no users in db")
	}
	res := make([]User, 0, len(usersInDB)-1)
	for _, v := range usersInDB {
		if v.ID != UserID(userId) {
			res = append(res, v)
		}
	}
	return res, nil
}

func InitStoreData(db commons.Datastore) error {
	sampleUsers := UsersInDB{
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

	userRoles := UserRoles{
		UserID("client_user"): []authz.Role{authz.RoleViewer},
	}

	// sample rbac state, with just one role viewer with permission to read all users
	rbac := authz.RbacInDB{
		authz.RoleViewer: authz.AccessRights{
			Role:        authz.RoleViewer,
			Resource:    authz.ResourceUser,
			Permissions: []authz.Permission{authz.PermissionRead},
			Conditions: authz.Conditions{
				authz.CondKeyResourceID: authz.ResourceIDAny,
			},
		},
	}
	db.Set("users", sampleUsers)
	db.Set("user_roles", userRoles)
	db.Set("rbac", rbac)
	//ignoring error handling in this case for this sample service
	return nil
}
