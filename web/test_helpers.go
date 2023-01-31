package web

import (
	"fmt"
	"github.com/MicBun/go-microservice-kubernetes/database"
	"github.com/MicBun/go-microservice-kubernetes/service"
	"io"
	"net/http"
	"net/http/httptest"
)

type TestFunc func(*service.Container)

func RunTest(testFunc TestFunc) {
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	if err := database.Migrate(db); err != nil {
		panic(err)
	}

	tx := db.Begin()
	defer tx.Rollback()
	c := service.New(tx)
	RegisterAPIRoutes(c)
	testFunc(c)
}

func MakeRequest(h http.Handler, method string, endpoint string, body io.Reader, headers ...map[string]string) (*httptest.ResponseRecorder, error) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(method, endpoint, body)
	if len(headers) > 0 {
		for key, val := range headers[0] {
			req.Header.Add(key, val)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("error building request %w", err)
	}
	h.ServeHTTP(w, req)
	return w, nil
}
