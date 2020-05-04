package controller

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"myWork/service"
)

type OrderController struct {
	Ctx iris.Context
	OrderService service.OrderService
	Session *sessions.Session
}

func (u *OrderController) GetQuery() mvc.Result {
	iris.New().Logger().Info("orderQuery")
	userId := u.Ctx.FormValue("user_id")
	order := u.OrderService.QueryOrderByUserId(userId)
	var resp  []interface{}
	for _, info := range order {
		resp = append(resp, info.TOrderToRespDesc())
	}
	return mvc.Response{
		Object: map[string]interface{}{
			"status": 0,
			"data"   :  resp,
		},
	}
}

func (u *OrderController) PostDelete() mvc.Result {
	iris.New().Logger().Info("orderDelete")

	orderId := u.Ctx.PostValue("order_id")
	startCity := u.Ctx.PostValue("start_city")
	endCity := u.Ctx.PostValue("end_city")
	seatNum, _:= u.Ctx.PostValueInt("seat_num")
	trainId := u.Ctx.PostValue("train_id")
	seatKind, _ := u.Ctx.PostValueInt("seat_kind")

	str := "退票失败"
	if u.OrderService.ReturnOrder(startCity, endCity, trainId, orderId, seatNum, seatKind) {
		str = "退票成功"
	}
	return mvc.Response{
		Object: map[string]interface{}{
			"status": 0,
			"data"   : str,
		},
	}
}
