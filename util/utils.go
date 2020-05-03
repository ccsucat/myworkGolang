package utils



import (
	"github.com/sirupsen/logrus"
	"reflect"
	"fmt"
	"github.com/kataras/iris"
	"time"

)

//根据Json格式设置obj对象
func SetObjByJson(obj interface{}, data map[string]interface{}) error {
	for key, value := range data {
		if err := setField(obj, key, value); err != nil {
			logrus.Error("SetObjByJson set field fail.")
			return err
		}
	}
	return nil
}

//设置结构体中的变量
func setField(obj interface{}, name string, value interface{}) error {
	structData := reflect.TypeOf(obj).Elem()
	fieldValue, result := structData.FieldByName(name)
	if !result {
		logrus.Error("No such field ", name)
		return fmt.Errorf("No such field %s", name)
	}

	//结构体中变量的类型
	fieldType := fieldValue.Type
	//参数的值
	val := reflect.ValueOf(value)
	//参数的类型
	valTypeStr := val.Type().String()
	//结构体中变量的类型
	fieldTypeStr := fieldType.String()
	//float64 to int
	if valTypeStr == "float64" && fieldTypeStr == "int" {
		val = val.Convert(fieldType)
	}

	//类型必须匹配
	if fieldType != val.Type() {
		return fmt.Errorf("value type %s didn't match obj field type %s ", valTypeStr, fieldTypeStr)
	}

	//fieldValue.Set(val)

	return nil
}

func StatusToNum(status int64) int {
	num := 0
	for i := 0; i < 60; i++ {
		if (status >> i & 1) == 0 {
			num++
		}
	}
	return num
}

func LogInfo(app *iris.Application, v ...interface{}) {
	app.Logger().Info(v)
}

func LogError(app *iris.Application, v ...interface{}) {
	app.Logger().Error(v)
}

func LogDebug(app *iris.Application, v ...interface{}) {
	app.Logger().Debug(v)
}

/**
 * 格式化数据
 */

func FormatTrainDatatime(startTime, endTime time.Time) string {
	t := endTime.Sub(startTime)
	str := fmt.Sprintf("%v小时%v分钟", int(t.Hours() / 1), int(t.Minutes()) % 60)
	return str
}

func FormatDatetime(time time.Time) string {
	return time.Format("2006-01-02 03:04")
}

func GetSeatNum(status int64) int {
	for i := 0; i < 60; i++ {
		if (status >> i & 1) == 0 {
			return i + 1
		}
	}
	return -1
}

func GetTimeByString(timeStr string) time.Time {
	timeLayout := "2006-01-02 15:04:05"  //转化所需模板
	local, _ := time.LoadLocation("Local")    //获取时区
	time, _ := time.ParseInLocation(timeLayout, timeStr, local)
	return time
}

