package server

import (
	"context"
	"net/http"
	"strings"
	"user-service/authz"
	"user-service/errorx"
	"user-service/users"
)

type Method string
type Path string
type MiddlewareFlags struct {
	AuthN bool
	AuthZ authz.AccessRights
}

// This middlewareOpts map can very well be stored in db, and populated in memory/cache during app init.
// But, it is fine to hardcode here for this sample service.
var middlewareOpts = map[Method]map[Path]MiddlewareFlags{
	http.MethodGet: {
		basePath + "/users": {AuthN: true, AuthZ: authz.AccessRights{
			Role:        authz.RoleViewer,
			Resource:    authz.ResourceUser,
			Permissions: []authz.Permission{authz.PermissionRead},
			Conditions:  authz.Conditions{authz.CondKeyResourceID: authz.ResourceIDAny},
		}},
		basePath + "/token": {},
	},
}

func getApiMiddlewareOpts(method Method, reqUrl string) MiddlewareFlags {
	pathOpts := middlewareOpts[method]
	for path, opts := range pathOpts {
		if strings.HasPrefix(reqUrl, string(path)) {
			return opts
		}
	}
	return MiddlewareFlags{}
}

func (a *App) AuthenticationMiddleware(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		opts := getApiMiddlewareOpts(Method(r.Method), r.URL.Path)

		// If authentication is not needed for a request, skip the checks
		if !opts.AuthN {
			inner.ServeHTTP(w, r)
			return
		}

		bearerToken := r.Header.Get("Authorization")
		if bearerToken == "" {
			RespondWithData(w, r, http.StatusBadRequest, errorx.Error{Code: errorx.BadRequestData, Message: "No API Key found"})
			return
		}
		tokenArray := strings.SplitAfter(bearerToken, "Bearer")
		if len(tokenArray) != 2 || strings.TrimSpace(tokenArray[1]) == "" {
			RespondWithData(w, r, http.StatusBadRequest, errorx.Error{Code: errorx.BadRequestData, Message: "Malformed API Key"})
			return
		}

		token := tokenArray[1]
		id, err := a.authNService.ValidateToken(strings.TrimSpace(token))
		if err != nil {
			switch err.Error() {
			case string(errorx.InvalidToken):
				RespondWithData(w, r, http.StatusUnauthorized, errorx.Error{Code: errorx.InvalidToken, Message: "Invalid API Key"})
			default:
				RespondWithData(w, r, http.StatusInternalServerError, errorx.Error{Code: errorx.ServerError})
			}
			return
		}

		// If token is valid, we put the id into the req context
		r = r.Clone(context.WithValue(r.Context(), users.UserIdInReqCtx, id))
		inner.ServeHTTP(w, r)
	})
}

func (a *App) AuthorizationMiddleware(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		opts := getApiMiddlewareOpts(Method(r.Method), r.URL.Path)
		if opts.AuthZ.Role == "" {
			inner.ServeHTTP(w, r)
			return
		}

		userId, ok := r.Context().Value(users.UserIdInReqCtx).(string)
		if !ok {
			RespondWithData(w, r, http.StatusBadRequest, errorx.Error{Code: errorx.BadRequestData, Message: "No authenticated user found"})
			return
		}
		roles, ok := a.db.Get("user_roles").(users.UserRoles)
		if !ok {
			RespondWithData(w, r, http.StatusInternalServerError, errorx.Error{Code: errorx.ServerError, Message: "No user_roles found to be matched"})
			return
		}
		role, ok := roles[users.UserID(userId)]
		if !ok {
			RespondWithData(w, r, http.StatusForbidden, errorx.Error{Code: errorx.AccessDenied, Message: "Forbidden. Insufficient Permissions"})
			return
		}
		// Since we know that we are dealing with just one role, we will take a shortcut for this sample service and get the role as role[0].
		// A similar shortcut we will take for Permissions as well, because we are dealing with just one permission.
		// In production scenario, these checks will be a bit more involved.
		authorized, err := a.authZService.IsAuthorized(string(role[0]), string(opts.AuthZ.Resource), string(opts.AuthZ.Permissions[0]), opts.AuthZ.Conditions)
		if err != nil {
			switch err.Error() {
			case string(errorx.AccessDenied):
				RespondWithData(w, r, http.StatusForbidden, errorx.Error{Code: errorx.AccessDenied, Message: "Forbidden. Insufficient Permissions"})
			default:
				RespondWithData(w, r, http.StatusInternalServerError, errorx.Error{Code: errorx.ServerError})
			}
			return
		}
		if !authorized {
			RespondWithData(w, r, http.StatusForbidden, errorx.Error{Code: errorx.AccessDenied, Message: "Forbidden. Insufficient Permissions"})
			return
		}

		// If authorization check is successful, we continue normally. There is nothing else to be put in the req context.
		inner.ServeHTTP(w, r)
	})
}
