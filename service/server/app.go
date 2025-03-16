package server

import (
	"context"
	"log"
	"sync"
	"user-service/authn"
	"user-service/authz"
	"user-service/commons"
	"user-service/config"
	"user-service/datastore"
)

type App struct {
	ctx          context.Context
	waitgroup    *sync.WaitGroup
	config       *config.Config
	db           commons.Datastore
	authNService commons.Authenticator
	authZService commons.Authorizer
}

func getApp(ctx context.Context) *App {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal("error initializing app config")
	}

	store := datastore.InitStore()
	authZSvc := authz.InitService(store)
	authNSvc := authn.InitService()
	authNSvc.Cfg = cfg

	// Initialize App
	a := App{
		ctx:          ctx,
		waitgroup:    &sync.WaitGroup{},
		config:       cfg,
		db:           store,
		authNService: authNSvc,
		authZService: authZSvc,
	}
	return &a
}
