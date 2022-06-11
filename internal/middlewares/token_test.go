package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/keystop/yaDiploma.git/internal/models"
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

func TestCheckAuthorization(t *testing.T) {
	uRM := new(UsersRepoMock)

	testValue := map[string]struct {
		token     string
		resStatus int
	}{
		"test1": {
			token:     "",
			resStatus: http.StatusUnauthorized,
		},
		"test2": {
			token:     "Bearer asdflkajfhkajdf",
			resStatus: http.StatusAccepted,
		},
		"test3": {
			token:     "Bearer ",
			resStatus: http.StatusUnauthorized,
		},
	}

	r := httptest.NewRequest("GET", "/", strings.NewReader(""))

	ur := new(models.User)

	tHandler := new(testHandler)
	funcHandler := CheckAuthorization(uRM)
	handler := funcHandler(tHandler)
	var w *httptest.ResponseRecorder

	for _, v := range testValue {
		w = httptest.NewRecorder()
		r.Header.Set("Authorization", v.token)
		if len(v.token) > 7 {
			ur.Token = v.token[7:]
		}
		uRM.On("Get", r.Context(), ur).Return(true, nil)
		handler.ServeHTTP(w, r)
		res := w.Result()
		defer res.Body.Close()
		assert.Equal(t, v.resStatus, res.StatusCode, "Не верный код ответа GET")
	}

}
