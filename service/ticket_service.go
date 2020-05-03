package service

import (
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris"
	"myWork/model"
	utils "myWork/util"
	"reflect"
	"sort"
)

type TicketService interface {
	GetTicketByCity(startCity, endCity string) []model.Travel
	BuyTicket(startCity, endCity, trainId string, seatNum, seatKind int) (model.TravelList, bool)
	AddOrder(order model.TOrder) bool
	UpdateTicket(list model.TravelList, seatNum, seatKind int)
}

func NewTicketService(engine *xorm.Engine) TicketService {
	return &ticketService{
		Engine: engine,
	}
}

type ticketService struct {
	Engine *xorm.Engine
}

func (u *ticketService) UpdateTicket(list model.TravelList, seatNum, seatKind int) {
	for _, info := range list {
		if seatKind == 0 {
			info.ZeroStatus |= (int64(1) << seatNum - 1)
		} else if seatKind == 1 {
			info.FirstStatus |= (int64(1) << seatNum - 1)
		} else {
			info.SecondStatus |= (int64(1) << seatNum - 1)
		}
		u.Engine.Where("train_id = ? and start_city = ?", info.TrainId, info.StartCity).Update(info)
	}
}

func (u *ticketService) BuyTicket(startCity, endCity, trainId string, seatNum, seatKind int) (model.TravelList, bool) {
	resp := model.TravelList{}
	value := model.Travel{}
	status := int64(0)
	for ;startCity != endCity; {
		u.Engine.Where("start_city = ? and train_id = ?", startCity, trainId).Get(&value)
		iris.New().Logger().Info(startCity, value.StartCity, value.EndCity)
		if reflect.DeepEqual(value, model.Travel{}) {
			return nil, false
		}
		resp = append(resp, value)
		if seatKind == 0 {
			status |= value.ZeroStatus
		} else if seatKind == 1 {
			status |= value.FirstStatus
		} else {
			status |= value.SecondStatus
		}
		startCity = value.EndCity
		value = model.Travel{}
	}
	if (status >> (seatNum - 1) & 1) == 1 {
		return nil, false
	}
	return resp, true
}

func (u *ticketService) GetTicketByCity(startCity, endCity string) []model.Travel {
	travel := []model.Travel{}
	u.Engine.Where("start_city = ?", startCity).Find(&travel)

	stack := utils.NewStack()
	for _, info := range travel {
		stack.Push(info)
	}
	num := 0
	resp := model.TravelList{}
	for ;stack.Len() != 0; {
		num++
		info, _ := stack.Top().(model.Travel)
		stack.Pop()
		if (info.EndCity == endCity) {
			resp = append(resp, info)
			continue
		}
		value := model.Travel{}
		iris.New().Logger().Info(info)
		u.Engine.Where("start_city = ? and train_id = ?", info.EndCity, info.TrainId).Get(&value)
		if reflect.DeepEqual(value, model.Travel{}) {
			continue
		}
		iris.New().Logger().Info(num, value.EndCity, info.SecondPrice, value.SecondPrice)
		info.EndCity = value.EndCity
		info.EndTime = value.EndTime
		info.ZeroPrice += value.ZeroPrice
		info.FirstPrice += value.FirstPrice
		info.SecondPrice += value.SecondPrice
		info.ZeroStatus |= value.ZeroStatus
		info.FirstStatus |= value.FirstStatus
		info.SecondStatus |= value.SecondStatus
		stack.Push(info)
	}
	sort.Sort(resp)
	return resp
}

func (u *ticketService) AddOrder(order model.TOrder) bool {
	_, err := u.Engine.Insert(order)
	if err != nil {
		iris.New().Logger().Info(err.Error())
	}
	return err == nil
}