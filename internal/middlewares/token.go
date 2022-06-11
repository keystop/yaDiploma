package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/keystop/yaDiploma/internal/models"
	"github.com/keystop/yaDiploma/pkg/logger"
)

func CheckAuthorization(ur models.UsersRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenName := "bearer "
			t := r.Header.Get("Authorization")
			key := ""
			if strings.HasPrefix(strings.ToLower(t), tokenName) {
				key = t[len(tokenName):]
			}

			if len(key) == 0 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			u := new(models.User)
			u.Token = key

			finded, err := ur.Get(r.Context(), u)
			if err != nil {
				logger.Info(http.StatusInternalServerError)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if !finded {
				logger.Info(http.StatusUnauthorized)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.Header().Add("Authorization", t)
			ctx := context.WithValue(r.Context(), models.UKeyName, u.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
