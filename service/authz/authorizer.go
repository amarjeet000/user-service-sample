// Package authz handles the authorization based on RBAC.
package authz

import (
	"log"
	"slices"
	"user-service/commons"
	"user-service/errorx"
)

type Role string
type Resource string
type Permission string
type CondKey string
type Conditions map[CondKey]interface{}

// AccessRights specifies a schema for RBAC policy.
// An AccessRights policy indicates which Role has what kind of Permission(s) on what Resource under what Conditions.
type AccessRights struct {
	Role        Role
	Resource    Resource
	Permissions []Permission
	Conditions  Conditions
}

/*
Note about RBAC (or ABAC - attribute-based-access-control) data structure:

The AccessRights data structure has been kept relatively simple, except the Conditions part that offers some flexibility.
Usually, a rich data structure is required depending up the complexity of the authorization policies.
In my experience, such a policy schema depends heavily upon the usecase.
It is also possible to keep both a rich policy schema as well as a simple RBAC schema side-by-side or in control of different services.
These two types of schema work together in deciding the final authorization for a user.

For this sample user-service, we will continue with a relatively simple AccessRights data structure.
*/

// Some hardcoded states. Ideally, they should be kept in a datastore.
const (
	RoleViewer Role = "viewer"

	ResourceUser Resource = "user"

	PermissionRead Permission = "read"

	CondKeyResourceID CondKey = "resource_id"

	ResourceIDAny = "*"
)

type RbacInDB map[Role]AccessRights

// Service implements Authorizer interface
type Service struct {
	store commons.Datastore
}

// Sample initialization
func InitService(db commons.Datastore) *Service {
	return &Service{store: db}
}

// IsRoleAuthorized checks if a Role has the requested Permission on a requested Resource under the requested Conditions.
func IsRoleAuthorized(db commons.Datastore, role string, resource string, permission string, conditions interface{}) (bool, error) {
	rbacInDB, ok := db.Get("rbac").(RbacInDB)
	if !ok {
		log.Println("DEBUG: rbac not found in store")
		return false, errorx.Error{Code: errorx.AccessDenied}
	}
	// get the access-rights (ar) of the supplied role
	aRights, ok := rbacInDB[Role(role)]
	if !ok {
		log.Println("DEBUG: role not found")
		return false, errorx.Error{Code: errorx.AccessDenied}
	}
	expectedConds, ok := conditions.(Conditions)
	if !ok {
		log.Println("DEBUG: malformed rbac conditions")
		return false, errorx.Error{Code: errorx.AccessDenied}
	}
	// Ensure that the expected resource matches the one in the access-rights retrieved from db.
	if aRights.Resource != Resource(resource) {
		log.Println("DEBUG: rbac resource mismatch")
		return false, errorx.Error{Code: errorx.AccessDenied}
	}
	// Ensure that the expected permission is part of the permissions list retrieved from db.
	if !slices.Contains(aRights.Permissions, Permission(permission)) {
		log.Println("DEBUG: requested permission not allowed")
		return false, errorx.Error{Code: errorx.AccessDenied}
	}
	// Ensure that the expected conditions match the one in the access-rights retrieved from db.
	// Since we know that we are dealing with just one condition key in this sample service, we will just check the match for that key.
	if aRights.Conditions[CondKeyResourceID] != expectedConds[CondKeyResourceID] {
		log.Println("DEBUG: rbac conditions mismatch")
		return false, errorx.Error{Code: errorx.AccessDenied}
	}

	return true, nil
}

func (s *Service) IsAuthorized(role string, resource string, permission string, conditions interface{}) (bool, error) {
	return IsRoleAuthorized(s.store, role, resource, permission, conditions)
}
