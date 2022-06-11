package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/keystop/yaDiploma.git/internal/models"
	"github.com/keystop/yaDiploma.git/pkg/logger"
)

func getUserFromRequest(r *http.Request) (*models.User, bool) {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.Info("Ошибка обработки запроса", err)
		return nil, false
	}
	user := new(models.User)

	err = json.Unmarshal(b, user)
	if err != nil {
		logger.Info("Ошибка Unmarshal", err)
		return nil, false
	}
	return user, true
}

func HandlerRegistration(ur models.UsersRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Обработка запроса регистрации")
		ctx := r.Context()

		user, ok := getUserFromRequest(r)
		if !ok {
			logger.Info(http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		finded, err := ur.Get(ctx, user)
		if err != nil {
			logger.Info(http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if finded {
			logger.Info(http.StatusConflict)
			w.WriteHeader(http.StatusConflict)
			return
		}

		if err := ur.Add(ctx, user); err != nil {
			logger.Info(http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Authorization", "Bearer "+user.Token)
		w.WriteHeader(http.StatusOK)
		logger.Info(http.StatusOK)
	}
}

func HandlerLogin(ur models.UsersRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Обработка запроса входа")
		ctx := r.Context()

		user, ok := getUserFromRequest(r)
		if !ok {
			logger.Info(http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logPass := user.Password

		finded, err := ur.Get(ctx, user)
		if err != nil {
			logger.Info(http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !finded || user.Password != logPass {
			logger.Info(http.StatusUnauthorized)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.Header().Add("Authorization", "Bearer "+user.Token)
		w.WriteHeader(http.StatusOK)
		logger.Info(http.StatusOK)
	}
}
