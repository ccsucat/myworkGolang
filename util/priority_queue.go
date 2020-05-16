package utils

import (
	"time"
)

type XTravel struct {
	TravelId      int
	City 		  string
	StartTime     time.Time
	EndTime		  time.Time
	TrainId    	  string
	ZeroPrice     float32
	FirstPrice    float32
	SecondPrice   float32
	ZeroStatus    int64
	FirstStatus   int64
	SecondStatus  int64
	Duration      int64
	ChangeTime    int
}
type T []XTravel

func (t *T) Len() int {
	return len(*t) //
}

func (t *T) Less(i, j int) bool {
	resp := false
	A := (*t)[i]
	B := (*t)[j]
	if A.ChangeTime == B.ChangeTime {
			return  A.Duration > B.Duration
	} else {
		return A.ChangeTime > B.ChangeTime
	}
	return resp
}

func (t *T) Swap(i, j int) {
	(*t)[i], (*t)[j] = (*t)[j], (*t)[i]
}

func (t *T) Push(x interface{}) {
	*t = append(*t, x.(XTravel))
}
func (t *T) Pop() interface{} {
	n := len(*t)
	x := (*t)[n-1]
	*t = (*t)[:n-1]
	return x
}
//func main() {
//	student := &Stu{{"Amy", 21}, {"Dav", 15}, {"Spo", 22}, {"Reb", 11}}
//	heap.Init(student)
//	one := stu{"hund", 9}
//	heap.Push(student, one)
//	for student.Len() > 0 {
//		fmt.Printf("%v\n", heap.Pop(student))
//	}
//}