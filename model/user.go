package model

import (
	"myWork/util"
	"time"
)

/**
 * 用户信息结构体,用于生成用户信息表
 */

type UserRegister struct {
	Id           string    `json:"id"`
	UserName     string    `json:"user_name"`
	Mobile       string    `json:"mobile"`
	Pwd          string    `json:"pwd"`
}

type User struct {
	Id           string     `xorm:"pk" json:"id"`        //主键 用户ID
	UserName     string    `xorm:"varchar(12)" json:"user_name"`  //用户名称
	RegisterTime time.Time `json:"register_time"`                //用户注册时间
	Mobile       string    `xorm:"varchar(11)" json:"mobile"`    //用户的移动手机号
	Balance      float64     `json:"balance"`                      //用户的账户余额（简单起见，使用int类型）
	Pwd          string    `json:"password"`                     //用户的账户密码
}

/**
 * 将数据库查询出来的结果进行格式组装成request请求需要的json字段格式
 */
func (user *User) UserToRespDesc() interface{} {
	respInfo := map[string]interface{}{
		"id":           user.Id,
		"user_name":    user.UserName,
		"registe_time": utils.FormatDatetime(user.RegisterTime),
		"mobile":       user.Mobile,
		"balance":      user.Balance,
	}
	return respInfo
}
