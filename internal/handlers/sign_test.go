package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/keystop/yaDiploma.git/internal/models"
	"github.com/keystop/yaDiploma.git/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type UsersRepoMock struct {
	mock.Mock
}

func (m *UsersRepoMock) Get(ctx context.Context, u *models.User) (bool, error) {
	args := m.Called(ctx, u)
	return args.Bool(0), args.Error(1)
}

func (m *UsersRepoMock) Add(ctx context.Context, u *models.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *UsersRepoMock) Del(ctx context.Context, u *models.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

type testHandler struct{}

func (t *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}

func TestHandlerLogin(t *testing.T) {
	testValue := map[string]struct {
		body       string
		userFinded bool
		findErr    error
		resStatus  int
	}{
		"test1": {
			userFinded: false,
			findErr:    nil,
			resStatus:  http.StatusBadRequest,
		},
		"test2": {
			body:       `{"login":"123", "password":"123"}`,
			userFinded: false,
			findErr:    errors.New(""),
			resStatus:  http.StatusInternalServerError,
		},
		"test3": {
			body:       `{"login":"123", "password":"123"}`,
			userFinded: false,
			findErr:    nil,
			resStatus:  http.StatusUnauthorized,
		},
		"test4": {
			body:       `{"login":"123", "password":"123"}`,
			userFinded: true,
			findErr:    errors.New("123"),
			resStatus:  http.StatusInternalServerError,
		},
		"test5": {
			body:       `{"login":"123", "password":"123"}`,
			userFinded: true,
			findErr:    nil,
			resStatus:  http.StatusOK,
		},
	}

	ur := &models.User{
		Login:    "123",
		Password: "123",
	}

	var w *httptest.ResponseRecorder

	for k, v := range testValue {

		urm := new(UsersRepoMock)
		tHandler := HandlerLogin(urm)
		r := httptest.NewRequest("GET", "/", strings.NewReader(v.body))
		w = httptest.NewRecorder()
		urm.On("Get", r.Context(), ur).Return(v.userFinded, v.findErr)
		tHandler.ServeHTTP(w, r)
		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, v.resStatus, res.StatusCode, k, "Не верный код ответа GET")
	}
}

func TestHandlerRegistration(t *testing.T) {

	testValue := map[string]struct {
		body       string
		userFinded bool
		findErr    error
		addErr     error
		resStatus  int
	}{
		"test1": {
			userFinded: false,
			addErr:     nil,
			findErr:    errors.New(""),
			resStatus:  http.StatusBadRequest,
		},
		"test2": {
			body:       `{"login":"123", "password":"123"}`,
			userFinded: false,
			findErr:    errors.New(""),
			resStatus:  http.StatusInternalServerError,
		},
		"test3": {
			body:       `{"login":"123", "password":"123"}`,
			userFinded: true,
			findErr:    errors.New("123"),
			resStatus:  http.StatusInternalServerError,
		},
		"test4": {
			body:       `{"login":"123", "password":"123"}`,
			userFinded: true,
			addErr:     errors.New("123"),
			findErr:    nil,
			resStatus:  http.StatusConflict,
		},
		"test5": {
			body:       `{"login":"123", "password":"123"}`,
			userFinded: false,
			addErr:     nil,
			findErr:    nil,
			resStatus:  http.StatusOK,
		},
	}

	ur := &models.User{
		Login:    "123",
		Password: "123",
	}

	var w *httptest.ResponseRecorder

	for k, v := range testValue {

		urm := new(UsersRepoMock)
		tHandler := HandlerRegistration(urm)
		r := httptest.NewRequest("POST", "/", strings.NewReader(v.body))
		w = httptest.NewRecorder()
		urm.On("Get", r.Context(), ur).Return(v.userFinded, v.findErr)
		urm.On("Add", r.Context(), ur).Return(v.addErr)
		tHandler.ServeHTTP(w, r)
		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, v.resStatus, res.StatusCode, k, "Не верный код ответа GET")
	}
}

func init() {
	logger.NewLogs()
	defer logger.Close()
}
