package model

import "time"

type TOrder struct {
	OrderId       int        `json:"order_id"`
	TrainId       string     `json:"train_id"`
	UserId        string     `json:"user_id"`
	UserName      string     `json:"user_name"`
	StartCity     string    `json:"start_city"`
	EndCity       string     `json:"end_city"`
	StartTime     time.Time  `json:"start_time"`
	EndTime       time.Time  `json:"end_time"`
	Price         float32    `json:"price"`
	SeatKind      int        `json:"seat_kind"`
	SeatNum       int        `json:"seat_num"`
}
