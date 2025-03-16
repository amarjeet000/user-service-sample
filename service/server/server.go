package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"user-service/config"
	"user-service/users"

	"github.com/go-chi/chi/v5"
)

var timeout = time.Second * 30

type UserService struct {
	router  *chi.Mux
	configs *config.Config
	wg      *sync.WaitGroup
}

func (s *UserService) init(ctx context.Context) {

	a := getApp(ctx)
	s.router = router(a)
	s.configs = a.config
	s.wg = a.waitgroup

	// We do some db state init here, ignoring error handling in this case for this sample service.
	users.InitStoreData(a.db)
}

func GetService(ctx context.Context) *UserService {
	// Initialize service
	s := UserService{}
	s.init(ctx)
	return &s
}

func (s *UserService) StartServer(ctx context.Context, cancel context.CancelFunc) {
	srv := http.Server{
		Addr:         s.configs.Host + ":" + s.configs.Port,
		Handler:      s.router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("an error occured, exiting from HTTP server", err)
		}
	}()

	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-quitSignal

	cancel()
	log.Println("INFO: gracefully shutting down...")

	ctxWithTimeOut, _ := context.WithTimeout(context.Background(), timeout)
	// We start the server shutdown with the provided timeout
	if err := srv.Shutdown(ctxWithTimeOut); err != nil {
		log.Println("ERROR: error shutting down the server via srv.Shutdown: ", err)
	}

	s.wg.Wait()
	log.Println("INFO: server shut down gracefully")
}
