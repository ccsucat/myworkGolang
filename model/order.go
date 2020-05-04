package model

import (
	utils "myWork/util"
	"time"
)

type TOrder struct {
	OrderId       int        `json:"order_id"`
	TrainId       string     `json:"train_id"`
	UserId        string     `json:"user_id"`
	UserName      string     `json:"user_name"`
	StartCity     string     `json:"start_city"`
	EndCity       string     `json:"end_city"`
	StartTime     time.Time  `json:"start_time"`
	EndTime       time.Time  `json:"end_time"`
	OrderTime     time.Time  `json:"order_time"`
	Price         float32    `json:"price"`
	SeatKind      int        `json:"seat_kind"`
	SeatNum       int        `json:"seat_num"`
	IsFirst       int        `json:"is_first"`
}

func (t *TOrder) TOrderToRespDesc() interface{} {
	respInfo := map[string]interface{}{
		"order_id"      :   t.OrderId,
		"train_id"      :   t.TrainId,
		"user_id"       :   t.UserId,
		"user_name"     :   t.UserName,
		"start_city"    :   t.StartCity,
		"end_city"      :   t.EndCity,
		"start_time"    :   utils.FormatDatetime(t.StartTime),
		"end_time"      :   utils.FormatDatetime(t.EndTime),
		"travel_time"   :   utils.FormatTrainDatatime(t.StartTime, t.EndTime),
		"order_time"    :   utils.FormatDatetime(t.OrderTime),
		"seat_kind"     :   t.SeatKind,
		"seat_num"      :   t.SeatNum,
		"price"         :   t.Price,
		"is_first"      :   t.IsFirst,

	}
	return respInfo
}
