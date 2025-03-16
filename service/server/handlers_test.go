package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"user-service/config"
	"user-service/testutils"
	"user-service/users"

	"github.com/go-chi/chi/v5"
)

var (
	configs      *config.Config
	store        *testutils.TestStore
	testAuthNSvc *testutils.TestAuthNService
	testAuthZSvc *testutils.TestAuthZService
)

func TestMain(m *testing.M) {
	// Setup before the tests
	store = testutils.InitTestStore()
	testAuthNSvc = testutils.InitTestAuthNService()
	testAuthZSvc = testutils.InitTestAuthZService(store)
	cfg, err := config.GetConfig()
	if err != nil {
		log.Println("Error instantiating cofig")
	}
	configs = cfg

	// Run tests and exit
	exitCode := m.Run()
	os.Exit(exitCode)
}

func testRouter() *chi.Mux {
	app := &App{
		ctx:          context.Background(),
		db:           store,
		authNService: testAuthNSvc,
		authZService: testAuthZSvc,
	}
	return router(app)
}

func TestGetUsers(t *testing.T) {
	router := testRouter()

	token, err := testAuthNSvc.GenerateToken("client_user")
	if err != nil {
		log.Println("Error during GenerateT", err)
	}
	tn := fmt.Sprintf("Bearer %s", token)
	headers := []testutils.Header{
		{Name: "Authorization", Value: tn},
	}

	w := testutils.MakeGetRequestWithHeaders(router, "/api/users", headers, []byte{})
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp []users.User
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		log.Println("Error processing resp", err)
	}
	if len(resp) != 2 { // and some other checks
		t.Errorf("unexpected response data, expected %v and %v", 2, resp)
	}

	// User with no permission
	token, err = testAuthNSvc.GenerateToken("bad_user")
	log.Println(token)
	if err != nil {
		log.Println("Error during GenerateT", err)
	}
	tn = fmt.Sprintf("Bearer %s", token)
	headers = []testutils.Header{
		{Name: "Authorization", Value: tn},
	}

	w = testutils.MakeGetRequestWithHeaders(router, "/api/users", headers, []byte{})
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}

	// Bad token case
	headers = []testutils.Header{
		{Name: "Authorization", Value: "bad_token"},
	}
	w = testutils.MakeGetRequestWithHeaders(router, "/api/users", headers, []byte{})
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}

	// No token cases
	headers = []testutils.Header{
		{Name: "Authorization", Value: ""},
	}
	w = testutils.MakeGetRequestWithHeaders(router, "/api/users", headers, []byte{})
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}

	headers = []testutils.Header{}
	w = testutils.MakeGetRequestWithHeaders(router, "/api/users", headers, []byte{})
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetToken(t *testing.T) {
	router := testRouter()
	headers := []testutils.Header{}
	w := testutils.MakeGetRequestWithHeaders(router, "/api/token", headers, []byte{})
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		log.Println("Error processing resp", err)
	}
	if resp["token"] == "" {
		t.Errorf("unexpected response data, expected %v and %v", 2, resp)
	}
}
