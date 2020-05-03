package controller

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"myWork/model"
	"myWork/service"
	"myWork/util"
	"strconv"
	"time"
)

//每一页最大的内容
const MaxLimit = 50

/**
 * 用户控制器结构体：用来实现处理用户模块的接口的请求，并返回给客户端
 */
type UserController struct {
	//上下文对象
	Ctx iris.Context
	//user service
	UserService service.UserService
	//session对象
	Session *sessions.Session
}


func(uc *UserController) PostRegister() mvc.Result {
	iris.New().Logger().Info("register")
	Id := uc.Ctx.PostValue("id")
	mobile := uc.Ctx.PostValue("mobile")

	userInfo := uc.UserService.GetUserInfoById(Id)
	iris.New().Logger().Info("--------", userInfo)
	if userInfo.Id != "" {
		return mvc.Response{
			Object: map[string]interface{} {
				"status": utils.REGISTERED_ID,
			},
		}
	}
	userInfo = uc.UserService.GetUserInfoByMobile(mobile)
	if userInfo.Id != "" {
		return mvc.Response{
			Object: map[string]interface{} {
				"status": utils.REGISTERED_MOBILE,
			},
		}
	}
	user := model.User{
		Id :            Id,
		UserName:       uc.Ctx.PostValue("user_name"),
		RegisterTime:   time.Now(),
		Mobile:         mobile,
		Balance:        10000,
		Pwd:            uc.Ctx.PostValue("pwd"),
	}
	if !uc.UserService.AddUser(user) {
		return mvc.Response{
			Object: map[string]interface{}{
				"status": utils.REGISTER_FAIL,
			},
		}
	}
	return mvc.Response{
		Object: map[string]interface{}{
			"status": utils.REGISTER_SUCCESS,
		},
	}
}

func (uc *UserController) PostLogin() mvc.Result {
	iris.New().Logger().Info("login")
	mobile := uc.Ctx.PostValue("mobile")
	pwd := uc.Ctx.PostValue("pwd")
	userInfo := uc.UserService.GetUserInfoByMobile(mobile)
	if userInfo.Pwd != pwd {
		return mvc.Response{
			Object: map[string]interface{}{
				"status": utils.LOGIN_FAIL,
			},
		}
	}
	return mvc.Response{
		Object: map[string]interface{}{
			"status": utils.LOGIN_SUCESS,
			"data": userInfo,
		},
	}
}

/**
 * 获取用户总数
 * 请求类型：Get
 * 请求Url：/v1/users/count
 */
func (uc *UserController) GetCount() mvc.Result {

	//用户总数
	total, err := uc.UserService.GetUserTotalCount()

	//请求出现错误
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				//"status": utils.RECODE_FAIL,
				"count":  0,
			},
		}
	}

	//正常情况的返回值
	return mvc.Response{
		Object: map[string]interface{}{
			//"status": utils.RECODE_OK,
			"count":  total,
		},
	}
}

func (uc *UserController)GetTemp() mvc.Result {

	return mvc.Response{
		Object: map[string]interface{}{
			//"status": utils.RECODE_OK,
			"data":  "123",
		},
	}
}

/**
 * 获取用户总数
 * 请求类型：Get
 * 请求Url：/v1/users/list
 */
func (uc *UserController) GetList() mvc.Result {

	offsetStr := uc.Ctx.FormValue("offset")
	limitStr := uc.Ctx.FormValue("limit")
	var offset int
	var limit int

	//判断offset和limit两个变量任意一个都不能为""
	if offsetStr == "" || limitStr == "" {

		return mvc.Response{
			Object: map[string]interface{}{
				//"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_USERLIST,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_USERLIST),
			},
		}
	}

	offset, err := strconv.Atoi(offsetStr)
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				//"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_USERLIST,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_USERLIST),
			},
		}
	}

	//做页数的限制检查
	if offset <= 0 {
		offset = 0
	}

	//做最大的限制
	if limit > MaxLimit {
		limit = MaxLimit
	}

	userList := uc.UserService.GetUserList(offset, limit)

	if len(userList) == 0 {
		return mvc.Response{
			Object: map[string]interface{}{
				//"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_USERLIST,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_USERLIST),
			},
		}
	}

	//将查询到的用户数据进行转换成前端需要的内容
	var respList []interface{}
	for _, user := range userList {
		respList = append(respList, user.UserToRespDesc())
	}

	//返回用户列表
	return mvc.Response{
		Object: &respList,
	}
}
