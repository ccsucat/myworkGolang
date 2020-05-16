package service

import (
	"container/heap"
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
	GetUserByUserId(userId, userName string) bool
	GetOrderByTime(startTime, endTime, userId string) bool
	GetTicketByInfo(startCity, endCity string, seatKind int, duration int64, price float32) []model.Travel
}

func NewTicketService(engine *xorm.Engine) TicketService {
	return &ticketService{
		Engine: engine,
	}
}

type ticketService struct {
	Engine *xorm.Engine
}

func (u *ticketService) GetTicketByInfo(startCity, endCity string, seatKind int, duration int64, price float32) []model.Travel {

	iris.New().Logger().Info("----", startCity, endCity, seatKind, duration, price)
	var pre [2000] int
	var vis [2000] int
	record := map[int]model.Travel{}
	cnt := 20
	var dp [2000] int64
	var f [2000] int
	var ans []int
	for i := 0; i < 2000; i++ {
		dp[i] = (int64(1) << 60)
		pre[i] = 0
		vis[i] = 0
		f[i] = 100000
	}
	travel := []model.Travel{}
	u.Engine.Where("start_city = ?", startCity).Find(&travel)
	queue := &utils.T{}
	temp := &utils.T{}
	//iris.New().Logger().Info("---------", travel)
	heap.Init(queue)
	heap.Init(temp)
	for _, info := range travel {
		queue.Push(info.TravelToX())
		pre[info.TravelId] = info.TravelId
		f[info.TravelId] = 1
		record[info.TravelId] = info
	}
	for ;queue.Len() > 0; {
		sort.Sort(queue)
		XTravel := queue.Pop().(utils.XTravel)
		iris.New().Logger().Info("----", XTravel.ChangeTime, XTravel.Duration)
		if vis[XTravel.TravelId] == 1 {
			continue
		}
		vis[XTravel.TravelId] = 1
		if XTravel.City == endCity {
			ans = append(ans, XTravel.TravelId)
			continue
			cnt--
			if cnt == 0 {
				break
			}
		}

		travel = []model.Travel{}
		u.Engine.Where("start_city = ?", XTravel.City).Find(&travel)
		for _, info := range travel {
			if XTravel.TrainId != info.TrainId {
				if (info.StartTime.Unix()-XTravel.EndTime.Unix()) > duration || info.StartTime.Unix() < XTravel.EndTime.Unix()+180 {
					//iris.New().Logger().Info("++++++++", (info.StartTime.Unix() - XTravel.EndTime.Unix()), duration)
					continue
				}
			}
			NewTravel := utils.XTravel{
				TravelId:     info.TravelId,
				City:         info.EndCity,
				StartTime:    XTravel.StartTime,
				EndTime:      info.EndTime,
				TrainId:      info.TrainId,
				ZeroPrice:    XTravel.ZeroPrice + info.ZeroPrice,
				FirstPrice:   XTravel.FirstPrice + info.FirstPrice,
				SecondPrice:  XTravel.SecondPrice + info.SecondPrice,
				ZeroStatus:   XTravel.ZeroStatus | info.ZeroStatus,
				FirstStatus:  XTravel.FirstStatus | info.FirstStatus,
				SecondStatus: XTravel.SecondStatus | info.SecondStatus,
				Duration:     info.EndTime.Unix() - XTravel.StartTime.Unix(),
				ChangeTime:   XTravel.ChangeTime,
			}
			if XTravel.TrainId != info.TrainId {
				NewTravel.ChangeTime++
			}
			ok := false
			Max := (int64(1) << 60) - 1
			if seatKind % 2 == 1 {
				if NewTravel.ZeroStatus != Max && NewTravel.ZeroPrice <= price {
					ok = true
				}
			}
			if seatKind % 4 >= 2 {
				if NewTravel.FirstStatus != Max && NewTravel.FirstPrice <= price {
					ok = true
				}
			}
			if seatKind >= 4 {
				if NewTravel.SecondStatus != Max && NewTravel.SecondPrice <= price {
					ok = true
				}
			}
			if !ok {
				continue
			}

			if f[info.TravelId] < NewTravel.ChangeTime || (f[info.TravelId] == NewTravel.ChangeTime && dp[info.TravelId] <= NewTravel.Duration) {
				continue
			}
			f[info.TravelId] = NewTravel.ChangeTime
			dp[info.TravelId] = NewTravel.Duration
			pre[info.TravelId] = XTravel.TravelId
			record[info.TravelId] = info
			heap.Push(queue, NewTravel)
			//queue.Push(NewTravel)
		}
	}
	//sort.Sort(temp)
	//for ;temp.Len() > 0; {
	//	tempTravel := temp.Pop().(utils.XTravel)
	//	ans = append(ans, tempTravel.TravelId)
	//}
	iris.New().Logger().Info("-----", ans)
	resp := model.TravelList{}
	for _, Tindex := range ans {
		stack := utils.NewStack()
		index := Tindex
		stack.Push(index)
		for ;pre[index] != index; {
			index = pre[index]
			stack.Push(index)
		}
		pre := model.Travel{}
		for ;stack.Len() != 0; {
			id := stack.Top().(int)
			stack.Pop()
			if reflect.DeepEqual(pre, model.Travel{}) {
				pre = record[id]
			} else {
				if pre.TrainId == record[id].TrainId {
					value := record[id]
					info := pre
					info.EndCity = value.EndCity
					info.EndTime = value.EndTime
					info.ZeroPrice += value.ZeroPrice
					info.FirstPrice += value.FirstPrice
					info.SecondPrice += value.SecondPrice
					info.ZeroStatus |= value.ZeroStatus
					info.FirstStatus |= value.FirstStatus
					info.SecondStatus |= value.SecondStatus
					pre = info
				} else {
					resp = append(resp, pre)
					pre = record[id]
				}
			}
		}
		resp = append(resp, pre)
	}
	//sort.Sort(resp)
	return resp
}


func (u *ticketService) GetOrderByTime(startTimeStr, endTimeStr, userId string) bool {
	startTimeDate := utils.GetTimeByString(startTimeStr)
	endTimeDate := utils.GetTimeByString(endTimeStr)
	startTimeUnix := startTimeDate.Unix() - int64(600)
	endTimeUnix := endTimeDate.Unix() + int64((600))
	order := model.TOrder{}
	u.Engine.Where("user_id = ? and start_time_unix <= ? and end_time_unix >= ?", userId, startTimeUnix, endTimeUnix).Get(&order)
	if order.UserId != "" {
		return true
	}
	order = model.TOrder{}
	u.Engine.Where("user_id = ? and start_time_unix <= ? and start_time_unix >= ?", userId, endTimeUnix, startTimeUnix).Get(&order)
	if order.UserId != "" {
		return true
	}
	order = model.TOrder{}
	u.Engine.Where("user_id = ? and end_time_unix <= ? and end_time_unix >= ?", userId, endTimeUnix, startTimeUnix).Get(&order)
	if order.UserId != "" {
		return true
	}
	return false


}

func (u *ticketService) GetUserByUserId(userId, userName string) bool {
	user := model.User{}
	u.Engine.Where("id = ? and user_name = ?", userId, userName).Get(&user)
	if user.Id == "" {
		return false
	}
	return true
}


func (u *ticketService) UpdateTicket(list model.TravelList, seatNum, seatKind int) {
	for _, info := range list {
		if seatKind == 0 {
			info.ZeroStatus |= (int64(1) << (seatNum - 1))
		} else if seatKind == 1 {
			info.FirstStatus |= (int64(1) << (seatNum - 1))
		} else {
			info.SecondStatus |= (int64(1) << (seatNum - 1))
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