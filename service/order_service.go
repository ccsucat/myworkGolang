package service

import (
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris"
	"myWork/model"
	"reflect"
)

type OrderService interface {
	AddOrder(order model.TOrder) bool
	QueryOrderByUserId(userId string) []model.TOrder
	ReturnOrder(startCity, endCity, trainId, orderId  string, seatNum, seatKind int) bool
}


func NewOrderService(engine *xorm.Engine) OrderService {
	return &orderService{
		Engine: engine,
	}
}

type orderService struct {
	Engine *xorm.Engine
}

func (u *orderService) QueryOrderByUserId(userId string) []model.TOrder {
	order := []model.TOrder{}
	u.Engine.Where("user_id = ?", userId).Find(&order)
	return order
}

func (u *orderService) AddOrder(order model.TOrder) bool {
	_, err := u.Engine.Insert(order)
	if err != nil {
		iris.New().Logger().Info(err.Error())
	}
	return err == nil
}

func (u *orderService) ReturnOrder(startCity, endCity, trainId, orderId string, seatNum, seatKind int) bool {
	order := model.TOrder{}
	u.Engine.Where("order_id = ?", orderId).Delete(order)
	value := model.Travel{}
	for ;startCity != endCity; {
		u.Engine.Where("start_city = ? and train_id = ?", startCity, trainId).Get(&value)
		if reflect.DeepEqual(value, model.Travel{}) {
			iris.New().Logger().Info("------------error")
			return false
		}

		if seatKind == 0 {
			value.ZeroStatus = value.ZeroStatus ^ (int64(1) << (seatNum - 1))
		} else if seatKind == 1 {
			value.FirstStatus = value.FirstStatus ^ (int64(1) << (seatNum - 1))
		} else {
			value.SecondStatus = value.SecondStatus ^ (int64(1) << (seatNum - 1))
		}
		u.Engine.Where("start_city = ? and train_id = ?", startCity, trainId).Update(value)
		startCity = value.EndCity
		value = model.Travel{}
	}
	return  true
}

