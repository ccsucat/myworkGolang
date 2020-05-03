package service

import (
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris"
	"myWork/model"
)

type OrderService interface {

}


func NewOrderService(engine *xorm.Engine) OrderService {
	return &orderService{
		Engine: engine,
	}
}

type orderService struct {
	Engine *xorm.Engine
}

func (u *orderService) AddOrder(order model.TOrder) bool {
	_, err := u.Engine.Insert(order)
	if err != nil {
		iris.New().Logger().Info(err.Error())
	}
	return err == nil
}

