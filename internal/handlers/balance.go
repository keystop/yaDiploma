package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/keystop/yaDiploma.git/internal/models"
	"github.com/keystop/yaDiploma.git/pkg/logger"
	"github.com/keystop/yaDiploma.git/pkg/luhn"
)

func HandlerGetUserBalance(br models.BalanceRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Запрос баланса пользователя")
		ctx := r.Context()
		userID := ctx.Value(models.UKeyName).(int)

		cb, err := br.Get(ctx, userID)
		if err != nil {
			logger.Info(http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		res, err := json.Marshal(cb)
		if err != nil {
			logger.Info("Ошибка маршализации", res)
			logger.Info(http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}

func HandlerGetUserWithdrawals(br models.BalanceRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Запрос списаний пользователя")
		ctx := r.Context()
		userID := ctx.Value(models.UKeyName).(int)

		cb, err := br.GetAll(ctx, userID)
		if err != nil {
			logger.Info(http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		res, err := json.Marshal(&cb)
		if err != nil {
			logger.Info("Ошибка маршализации", res)
			logger.Info(http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}

func HandlerGetUserWithdraw(br models.BalanceRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Запрос списания средств пользователем")
		ctx := r.Context()
		userID := ctx.Value(models.UKeyName).(int)

		b, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Info(http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		bo := new(models.BalanceOut)

		err = json.Unmarshal(b, &bo)
		if err != nil {
			logger.Info(http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if ok := luhn.CheckString(bo.OrderID); !ok || len(bo.OrderID) == 0 {
			logger.Info(http.StatusUnprocessableEntity)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		cb, err := br.Get(ctx, userID)
		if err != nil {
			logger.Info(http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if cb.CurBalance < bo.Sum {
			logger.Info(http.StatusPaymentRequired)
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}

		bl := &models.Balance{
			UserID:  userID,
			OrderID: bo.OrderID,
			SumOut:  bo.Sum,
		}
		if err = br.Add(ctx, bl); err != nil {
			logger.Info(http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logger.Info(http.StatusOK)
		w.WriteHeader(http.StatusOK)
	}
}
