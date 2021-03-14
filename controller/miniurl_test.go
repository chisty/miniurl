package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/chisty/miniurl/model"
	"github.com/gorilla/mux"
)

type mockURLService struct{}

func (msvc *mockURLService) Get(id string) (*model.MiniURL, error) {
	return createMiniUrl(id, "http://mock-test-url"), nil
}

func (msvc *mockURLService) Save(data *model.MiniURL) (*model.MiniURL, error) {
	return createMiniUrl(data.ID, data.URL), nil
}

type mockRedisWithoutValue struct{}

func (r *mockRedisWithoutValue) Set(key string, value *model.MiniURL) error {
	return nil
}

func (r *mockRedisWithoutValue) Get(key string) (*model.MiniURL, error) {
	return nil, errors.New("item not found")
}

type mockRedisWithValue struct{}

func (r *mockRedisWithValue) Set(key string, value *model.MiniURL) error {
	return nil
}

func (r *mockRedisWithValue) Get(key string) (*model.MiniURL, error) {
	return createMiniUrl("testID", "http://mock-test-url"), nil
}

func createMiniUrl(id string, url string) *model.MiniURL {
	return &model.MiniURL{
		ID:  id,
		URL: url,
	}
}

func createHandlerWithoutCache() http.HandlerFunc {
	logger := log.New(os.Stdout, "miniurl-app", log.LstdFlags|log.Lshortfile)
	ctrl := NewMiniURLCtrl(&mockURLService{}, &mockRedisWithoutValue{}, logger)
	return http.HandlerFunc(ctrl.Get)
}

func createHandlerWithCache() http.HandlerFunc {
	logger := log.New(os.Stdout, "miniurl-app", log.LstdFlags|log.Lshortfile)
	ctrl := NewMiniURLCtrl(&mockURLService{}, &mockRedisWithoutValue{}, logger)
	return http.HandlerFunc(ctrl.Get)
}

func TestGetInvalidRequest(t *testing.T) {
	handler := createHandlerWithoutCache()
	r, _ := http.NewRequest("GET", "/", nil)
	rw := httptest.NewRecorder()
	vars := map[string]string{
		"id": "  ",
	}

	r = mux.SetURLVars(r, vars)
	handler.ServeHTTP(rw, r)

	if rw.Code != http.StatusBadRequest {
		t.Error("Test failed. Expected http.StatusBadRequest")
	}

	if string(rw.Body.Bytes()) != "invalid request value in route" {
		t.Error("Test failed. Response body mismatched.")
	}
}

func TestGetServiceIfCacheMiss(t *testing.T) {
	handler := createHandlerWithoutCache()
	r, _ := http.NewRequest("GET", "/", nil)
	rw := httptest.NewRecorder()
	vars := map[string]string{
		"id": "testID",
	}

	r = mux.SetURLVars(r, vars)
	handler.ServeHTTP(rw, r)
	if rw.Code != http.StatusTemporaryRedirect {
		t.Error("Test failed. Expected http.StatusTemporaryRedirect")
	}
}

func TestGetServiceWithCache(t *testing.T) {
	handler := createHandlerWithCache()
	r, _ := http.NewRequest("GET", "/", nil)
	rw := httptest.NewRecorder()
	vars := map[string]string{
		"id": "testID",
	}

	r = mux.SetURLVars(r, vars)
	handler.ServeHTTP(rw, r)
	if rw.Code != http.StatusTemporaryRedirect {
		t.Error("Test failed. Expected http.StatusTemporaryRedirect")
	}
}

func TestPostValidRequest(t *testing.T) {
	logger := log.New(os.Stdout, "miniurl-app", log.LstdFlags|log.Lshortfile)
	ctrl := NewMiniURLCtrl(&mockURLService{}, &mockRedisWithoutValue{}, logger)

	tests := []struct {
		input    []byte
		expected string
	}{
		{[]byte(`{"url":"http://www.gogle.com"}`), "http://www.gogle.com"},
		{[]byte(`{"url":"https://www.gogle.com"}`), "https://www.gogle.com"},
		{[]byte(`{"url":"http://asd"}`), "http://asd"},
		{[]byte(`{"url":"https://ap-southeast-1.console.aws.amazon.com/dynamodb/home?region=ap-southeast-1#tables:selected=eatigo;tab=items"}`), "https://ap-southeast-1.console.aws.amazon.com/dynamodb/home?region=ap-southeast-1#tables:selected=eatigo;tab=items"},
	}

	for _, test := range tests {
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(test.input))
		handler := http.HandlerFunc(ctrl.Save)

		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Test failed. Expeced status %d, found %d\n", http.StatusOK, resp.Code)
		}

		miniurl := model.MiniURL{}
		err := json.Unmarshal(resp.Body.Bytes(), &miniurl)
		if err != nil {
			t.Error("Test failed. Cannot decode response")
		}

		if miniurl.URL != test.expected {
			t.Errorf("Test failed. Expeced URL %s, found %s\n", test.expected, miniurl.URL)
		}
	}
}

func TestPostInvalidRequest(t *testing.T) {
	logger := log.New(os.Stdout, "miniurl-app", log.LstdFlags|log.Lshortfile)
	ctrl := NewMiniURLCtrl(&mockURLService{}, &mockRedisWithoutValue{}, logger)

	tests := []struct {
		input    []byte
		expected string
	}{
		{[]byte(`{"url":"http//www.gogle.com"}`), "http//www.gogle.com"},
		{[]byte(`{"url":"_www.gogle.com"}`), "_www.gogle.com"},
		{[]byte(`{"url":"123:asd"}`), "123:asd"},
		{[]byte(`{"url":"SimpleText"}`), "SimpleText"},
	}

	for _, test := range tests {
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(test.input))
		handler := http.HandlerFunc(ctrl.Save)

		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)

		if resp.Code != http.StatusBadRequest {
			t.Errorf("Test failed. Expeced status %d, found %d, value=%s\n", http.StatusBadRequest, resp.Code, test.expected)
		}
	}
}
