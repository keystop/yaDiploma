package accrual

import (
	"context"
	"encoding/json"

	"github.com/keystop/yaDiploma/internal/config"
	"github.com/keystop/yaDiploma/internal/models"
	"github.com/keystop/yaDiploma/pkg/client"
	"github.com/keystop/yaDiploma/pkg/logger"
	"github.com/keystop/yaDiploma/pkg/utility"
	"github.com/keystop/yaDiploma/pkg/workers"
)

// servAddress string

type ordersDB struct {
	or models.OrdersRepo
	br models.BalanceRepo
}

func (o *ordersDB) getOrders(ctx context.Context, st models.OrderStatus) []*models.Order {
	arrOr, err := o.or.GetAllStatus(ctx, st)
	if err != nil {
		logger.Info("Ошибка запроса не обработанных заказов", err)
	}
	return arrOr
}

func (o *ordersDB) updateOrderDB(ctx context.Context, or *models.Order, oOut *models.OrderFromAccrual) {
	if or.Status != oOut.Status {
		or.OrderID = oOut.OrderID
		or.Status = oOut.Status
		or.Accrual = oOut.Accrual

		o.or.Update(ctx, or)

		if oOut.Status == models.OrderStatusProcessed {
			bl := new(models.Balance)
			bl.UserID = or.UserID
			bl.OrderID = or.OrderID
			bl.SumIn = or.Accrual
			err := o.br.Add(ctx, bl)
			if err != nil {
				logger.Info("Ошибка добавления баланса в лог", bl)
			}
		}
	}
}

// default:
// 	l.getNumForSurvey(ctx, models.OrderStatusNew)
// 	l.getNumForSurvey(ctx, models.OrderStatusRegistered)
// 	l.getNumForSurvey(ctx, models.OrderStatusProcessing)
// 	time.Sleep(1 * time.Second)
// }
// }

type Accrual struct {
	servAddress string
	odb         *ordersDB
	w           *workers.WorkersPool
}

func (a *Accrual) PutArr(ctx context.Context, arr []*models.Order) {
	for _, v := range arr {
		a.w.Put(a.getOrderFromAccrual(ctx, v))
	}
}

func (a *Accrual) Put(ctx context.Context, o *models.Order) {
	a.w.Put(a.getOrderFromAccrual(ctx, o))
}

func (a *Accrual) getOrdersForSurvey(ctx context.Context, st models.OrderStatus) {
	arr := a.odb.getOrders(ctx, st)
	a.PutArr(ctx, arr)
}

func (a *Accrual) GetOrdersForSurveyFromDB(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			a.getOrdersForSurvey(ctx, models.OrderStatusNew)
			a.getOrdersForSurvey(ctx, models.OrderStatusRegistered)
			a.getOrdersForSurvey(ctx, models.OrderStatusProcessing)
		}
	}
}

func (a *Accrual) getOrderFromAccrual(ctx context.Context, o *models.Order) func() {
	return func() {
		orderForUpdate := o
		oIn := new(models.OrderFromAccrual)
		b, ok := client.MakeRequest("GET", a.servAddress+"api/orders/"+orderForUpdate.OrderID, "", "", nil)
		if !ok {
			return
		}
		err := json.Unmarshal(b, oIn)
		if err != nil {
			logger.Info("Error", "Ошибка чтения ответа на запрос", err)
			return
		}

		a.odb.updateOrderDB(ctx, orderForUpdate, oIn)
	}
}

func NewSurveyAccrual(or models.OrdersRepo, br models.BalanceRepo, w *workers.WorkersPool) *Accrual {
	ordersDB := &ordersDB{
		or: or,
		br: br,
	}
	a := new(Accrual)
	a.odb = ordersDB
	a.servAddress = utility.IncludeTrailingBackSlash(config.Cfg.AccrualAddress())
	a.w = w

	return a
}
