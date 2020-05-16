package model

import (
	utils "myWork/util"
	"time"
)

type Travel struct {
	TravelId      int        `json:"travel_id"`
	TrainId       string     `json:"train_id"`
	StartCity     string     `json:"start_city"`
	EndCity       string     `json:"end_city"`
	StartTime     time.Time  `json:"start_time"`
	EndTime       time.Time  `json:"end_time"`
	ZeroPrice     float32    `json:"zero_price"`
	FirstPrice    float32    `json:"first_price"`
	SecondPrice   float32    `json:"second_price"`
	ZeroStatus    int64      `json:"zero_status"`
	FirstStatus   int64      `json:"first_status"`
	SecondStatus  int64      `json:"second_status"`
}

type TravelList []Travel

func (T TravelList) Len() int {
	return len(T)
}

func (T TravelList) Less(i, j int) bool {
	return T[i].StartTime.Unix() < T[j].StartTime.Unix()
}

func (T TravelList) Swap(i, j int) {
	T[i], T[j] = T[j], T[i]
}

type Ticket struct {
	TrainId       string
	StartCity     string
	EndCity       string
	StartTime     time.Time
	EndTime       time.Time
	durationTime  time.Time
	ZeroPrice     float32
	FirstPrice    float32
	SecondPrice   float32
}

func (t *Travel) TravelToX() utils.XTravel {
	return utils.XTravel{
		TravelId:     t.TravelId,
		City:         t.EndCity,
		StartTime:    t.StartTime,
		EndTime:      t.EndTime,
		TrainId:      t.TrainId,
		ZeroPrice:    t.ZeroPrice,
		FirstPrice:   t.FirstPrice,
		SecondPrice:  t.SecondPrice,
		ZeroStatus:   t.ZeroStatus,
		FirstStatus:  t.FirstStatus,
		SecondStatus: t.SecondStatus,
		Duration:     t.EndTime.Unix() - t.StartTime.Unix(),
		ChangeTime:   1,
	}
}

func (t *Travel) TravelToRespDesc() interface{} {
	respInfo := map[string]interface{}{
		"train_id"      :   t.TrainId,
		"start_city"    :   t.StartCity,
		"end_city"      :   t.EndCity,
		"start_time"    :   utils.FormatDatetime(t.StartTime),
		"end_time"      :   utils.FormatDatetime(t.EndTime),
		"travel_time"   :   utils.FormatTrainDatatime(t.StartTime, t.EndTime),
		"zero_price"    :   t.ZeroPrice,
		"first_price"   :   t.FirstPrice,
		"second_price"  :   t.SecondPrice,
		"zero_num"      :   utils.StatusToNum(t.ZeroStatus),
		"first_num"     :   utils.StatusToNum(t.FirstStatus),
		"second_num"    :   utils.StatusToNum(t.SecondStatus),
		"zero_seat"     :   utils.GetSeatNum(t.ZeroStatus),
		"first_seat"    :   utils.GetSeatNum(t.FirstStatus),
		"second_seat"   :   utils.GetSeatNum(t.SecondStatus),
	}
	return respInfo
}