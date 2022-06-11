package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

	"github.com/keystop/yaDiploma.git/internal/models"
	"github.com/keystop/yaDiploma.git/pkg/logger"
	"github.com/keystop/yaDiploma.git/pkg/luhn"
)

func HandlersNewOrder(or models.OrdersRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Обработка нового заказа")
		ctx := r.Context()
		userID := ctx.Value(models.UKeyName).(int)
		order := new(models.Order)
		order.UserID = userID

		b, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ordNum := string(b)
		if ok := luhn.CheckString(ordNum); len(ordNum) == 0 || !ok {
			logger.Info(ordNum, http.StatusUnprocessableEntity)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		order.OrderID = string(b)

		finded, err := or.Get(ctx, order)
		if err != nil {
			logger.Info(http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if finded {
			status := http.StatusOK
			if order.UserID != userID {
				status = http.StatusConflict
			}
			logger.Info(status)
			w.WriteHeader(status)
			return
		}
		if err := or.Add(ctx, order); err != nil {
			logger.Info(http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logger.Info(http.StatusAccepted)
		w.WriteHeader(http.StatusAccepted)
	}
}

func HandlersGetUserOrders(or models.OrdersRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Поиск заказов пользователя")
		ctx := r.Context()
		userID := ctx.Value(models.UKeyName).(int)

		arrOrders, err := or.GetAll(ctx, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				logger.Info(http.StatusNoContent)
				w.WriteHeader(http.StatusNoContent)
				return
			}
			logger.Info(http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		res, err := json.Marshal(&arrOrders)
		if err != nil {
			logger.Info("Ошибка маршализации", arrOrders)
			logger.Info(http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		logger.Info(http.StatusOK)
		w.WriteHeader(http.StatusOK)
		w.Write(res)

	}
}
