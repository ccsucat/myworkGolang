package controller

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"myWork/model"
	"myWork/service"
	utils "myWork/util"
	"strconv"
	"time"
)

type TicketController struct {
	Ctx iris.Context
	TicketService service.TicketService
	Session *sessions.Session
}

func (u *TicketController) GetQuery() mvc.Result {
	iris.New().Logger().Info("query")
	startCity := u.Ctx.FormValue("start_city")
	endCity := u.Ctx.FormValue("end_city")
	infos := u.TicketService.GetTicketByCity(startCity, endCity)

	var rep  []interface{}
	for _, info := range infos {
		rep = append(rep, info.TravelToRespDesc())
	}
	return mvc.Response{
		Object: map[string]interface{}{
			"status": 0,
			"data"   :  rep,
		},
	}
}

func (u *TicketController) GetQuery2() mvc.Result {
	iris.New().Logger().Info("moreQuery")
	startCity := u.Ctx.FormValue("start_city")
	endCity := u.Ctx.FormValue("end_city")
	priceStr := u.Ctx.FormValue("price")
	seatkindStr := u.Ctx.FormValue("seat_kind")
	durationStr := u.Ctx.FormValue("duration")
	price, _ := strconv.ParseFloat(priceStr,32)
	seatkind, _ := strconv.ParseInt(seatkindStr, 10, 32)
	duration, _ := strconv.ParseInt(durationStr, 10, 64)
	infos := u.TicketService.GetTicketByInfo(startCity, endCity, int(seatkind), duration * 60, float32(price))
	var rep  []interface{}
	for _, info := range infos {
		rep = append(rep, info.TravelToRespDesc())
	}
	return mvc.Response{
		Object: map[string]interface{}{
			"status": 0,
			"data"   :  rep,
		},
	}
}

func (u *TicketController)PostBuy() mvc.Result {
	iris.New().Logger().Info("buy")
	startCity := u.Ctx.PostValue("start_city")
	endCity := u.Ctx.PostValue("end_city")
	seatNum, _:= u.Ctx.PostValueInt("seat_num")
	trainId := u.Ctx.PostValue("train_id")
	seatKind, _ := u.Ctx.PostValueInt("seat_kind")
	userId := u.Ctx.PostValue("user_id")
	userName := u.Ctx.PostValue("user_name")
	price, _ := u.Ctx.PostValueFloat64("price")
	startTime := u.Ctx.PostValue("start_time")
	endTime := u.Ctx.PostValue("end_time")
	isFirst, _ := u.Ctx.PostValueInt("is_first")
	iris.New().Logger().Info(isFirst)
	if !u.TicketService.GetUserByUserId(userId, userName) {
		return mvc.Response{
			Object: map[string]interface{}{
				"status": 0,
				"data"   :  "没有该用户",
			},
		}
	}
	resp, ok := u.TicketService.BuyTicket(startCity, endCity, trainId, seatNum, seatKind)
	if !ok {
		return mvc.Response{
			Object: map[string]interface{}{
				"status": 0,
				"data"   :  "购买失败",
			},
		}
	}
	if u.TicketService.GetOrderByTime(startTime, endTime, userId) {
		return mvc.Response{
			Object: map[string]interface{}{
				"status": 0,
				"data"   :  "购买失败，与已有订单时间冲突",
			},
		}
	}
	order := model.TOrder{
		TrainId:   trainId,
		UserId:    userId,
		UserName:  userName,
		StartCity: startCity,
		EndCity:   endCity,
		StartTime: utils.GetTimeByString(startTime),
		EndTime:   utils.GetTimeByString(endTime),
		StartTimeUnix: utils.GetTimeByString(startTime).Unix(),
		EndTimeUnix:   utils.GetTimeByString(endTime).Unix(),
		OrderTime: time.Now(),
		Price:     float32(price),
		SeatKind:  seatKind,
		SeatNum:   seatNum,
		IsFirst:   isFirst,
	}
	ok = u.TicketService.AddOrder(order)
	if !ok {
		return mvc.Response{
			Object: map[string]interface{}{
				"status": 0,
				"data"   :  "购买失败",
			},
		}
	}
	u.TicketService.UpdateTicket(resp, seatNum, seatKind)
	return mvc.Response{
		Object: map[string]interface{}{
			"status": 1,
			"data"   :  "购买成功",
		},
	}

}