package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddle "github.com/go-chi/chi/v5/middleware"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

const basePath = "/api"

func router(app *App) *chi.Mux {
	r := chi.NewRouter()
	r.Use(
		// NOTE: A CORS middleware can be placed if there is need for it
		chimiddle.Logger,
		app.AuthenticationMiddleware,
		app.AuthorizationMiddleware,
	)
	routes := []Route{
		{
			// This is a helper open endpoint to get a token via an http request.
			// Usually, token issuance happens to a registered client.
			// During registration, the client gets issued a client_id and a client_secret.
			// At the time of registration, the client should also provide its public key to the server.
			// This public key is used to verify the signature of the client when the client initiates the process of token generation.
			// But, for this sample service, in the interest of time, this open endpoint is fine.
			Name:        "GetToken",
			Method:      "GET",
			Pattern:     basePath + "/token",
			HandlerFunc: app.GetToken,
		},
		{
			Name:        "GetUser",
			Method:      "GET",
			Pattern:     basePath + "/users",
			HandlerFunc: app.GetUsers,
		},
	}

	for _, v := range routes {
		r.MethodFunc(v.Method, v.Pattern, v.HandlerFunc)
	}

	return r
}
