package server

import (
	"net/http"
	"user-service/errorx"
	"user-service/users"
)

func (app *App) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Since the authentication and authorization have been taken care of at the middleware level,
	// the handler simply processes the request and respond appropriately.

	// Such a simple auth model is not always possible in a complex business scenario.
	// Often a fine-grained authorization check is needed at the handler level on top of the generic middlware checks.
	userId, ok := r.Context().Value(users.UserIdInReqCtx).(string)
	if !ok {
		RespondWithData(w, r, http.StatusBadRequest, errorx.Error{Code: errorx.BadRequestData, Message: "Bad API key"})
		return
	}

	usrs, err := users.FetchUsersFilterOne(app.db, userId)
	if err != nil {
		RespondWithData(w, r, http.StatusNoContent, errorx.Error{Code: errorx.NoContent, Message: "No users found"})
		return
	}

	RespondWithData(w, r, http.StatusOK, usrs)
}

func (app *App) GetToken(w http.ResponseWriter, r *http.Request) {
	token, err := app.authNService.GenerateToken("client_user")
	if err != nil {
		RespondWithData(w, r, http.StatusInternalServerError, errorx.Error{Code: errorx.ServerError, Message: "Could not generate API Key"})
		return
	}
	res := map[string]string{"token": token}
	RespondWithData(w, r, http.StatusOK, res)
}
