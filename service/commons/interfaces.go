// Package commons exposes facilities that can be consumed by more than one packages.
// Note that there are differences in opinion in the Go community on where to define interfaces in a Go project.
// For now, I will lean on the side of the readability in this sample service.
// But, it often depends upon the development team and the complexity of the project.
// I do not have a strong preference on this matter.
package commons

// Datastore exposes simple Get and Set methods
type Datastore interface {
	Get(key string) interface{}
	Set(key string, val interface{}) error
}

// Authenticator handles token generation and validation
type Authenticator interface {
	GenerateToken(id string) (string, error)
	ValidateToken(token string) (string, error)
}

// Authorizer exposes a method to check if a role has a required permission(s) on a resource under certain conditions.
type Authorizer interface {
	IsAuthorized(role string, resource string, permission string, conditions interface{}) (bool, error)
}
