package service

import (
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris"
	"myWork/model"
)

/**
 * 用户模块功能服务接口
 */
type UserService interface {

	GetUserInfoById(ID string) model.User

	GetUserInfoByMobile(mobile string) model.User

	AddUser(user model.User) bool

	//获取用户日增长统计数据
	GetUserDailyStatisCount(datetime string) int64
	//获取用户总数
	GetUserTotalCount() (int64, error)
	//用户列表
	GetUserList(offset, limit int) []*model.User
}

/**
 * 实例化用户服务结构实体对象
 */
func NewUserService(engine *xorm.Engine) UserService {
	return &userService{
		Engine: engine,
	}
}


/**
 * 用户服务实现结构体
 */
type userService struct {
	Engine *xorm.Engine
}

func (uc *userService) GetUserById(ID string) bool {
	return true
}

func (uc *userService) AddUser(user model.User) bool {
	_, err := uc.Engine.Insert(user)
	if err != nil {
		iris.New().Logger().Info(err.Error())
	}
	return err == nil
}

func (uc *userService)  GetUserInfoById(ID string) model.User {
	user:= model.User{}
	uc.Engine.Where("id = ?", ID).Get(&user)
	return user
}

func (uc *userService)  GetUserInfoByMobile(mobile string) model.User {
	var user model.User
	uc.Engine.Where("mobile = ?", mobile).Get(&user)
	return user
}



/**
 * 请求总的用户数量
 * 返回值：总用户数量
 */
func (uc *userService) GetUserTotalCount() (int64, error) {

	//查询del_flag 为0 的用户的总数量；del_flag:0 正常状态；del_flag:1 用户注销或者被删除
	count, err := uc.Engine.Where(" del_flag = ? ", 0).Count(new(model.User))
	if err != nil {
		panic(err.Error())
		return 0, err
	}
	//用户总数
	return count, nil
}

/**
 * 请求用户列表数据
 * offset：偏移数量
 * limit：一次请求获取的数据条数
 */
func (uc *userService) GetUserList(offset, limit int) []*model.User {

	var userList []*model.User

	err := uc.Engine.Where("del_flag = ?", 0).Limit(limit, offset).Find(&userList)
	if err != nil {
		iris.New().Logger().Error(err.Error())
		panic(err.Error())
		return nil
	}
	return userList
}

/**
 * 获取用户日增长统计结果
 * datetime：某一个特殊的日期
 */
func (us *userService) GetUserDailyStatisCount(datetime string) int64 {

	result, err := us.Engine.Count(new(model.User))
	if err != nil {
		panic(err.Error())
	}

	return result
}
