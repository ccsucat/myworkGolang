package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"myWork/controller"
	"myWork/datasource"
	"myWork/service"
	"time"
)

func main() {
	app := iris.New()

	configation(app)
	mvcHandle(app)

	//config := config2.InitConfig()
	addr := ":" + "7999"
	app.Run(
		iris.Addr(addr),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
		)

	//app.Get("/Get", func(context context.Context) {
	//	a := "number"
	//	b := 2
	//	context.JSON(iris.Map{
	//		a : b,
	//	})
	//	app.Logger().Info("hello world")
	//})
	//app.Run(iris.Addr(":8000"))

}

func mvcHandle(app *iris.Application) {
	sessManager := sessions.New(sessions.Config{
		Cookie:  "sessioncookie",
		Expires: 24 * time.Hour,
	})

	engine := datasource.NewMysqlEngine()
	userService := service.NewUserService(engine)
	user := mvc.New(app.Party("/user"))
	user.Register(
		userService,
		sessManager.Start,
		)
	user.Handle(new(controller.UserController))


	ticketService := service.NewTicketService(engine)
	ticket := mvc.New(app.Party("/ticket"))
	ticket.Register(
		ticketService,
		sessManager.Start,
		)
	ticket.Handle(new(controller.TicketController))



}

func configation(app *iris.Application) {

	//配置 字符编码
	app.Configure(iris.WithConfiguration(iris.Configuration{
		Charset: "UTF-8",
	}))

	//错误配置
	//未发现错误
	app.OnErrorCode(iris.StatusNotFound, func(context context.Context) {
		context.JSON(iris.Map{
			"errmsg": iris.StatusNotFound,
			"msg":    " not found ",
			"data":   iris.Map{},
		})
	})

	app.OnErrorCode(iris.StatusInternalServerError, func(context context.Context) {
		context.JSON(iris.Map{
			"errmsg": iris.StatusInternalServerError,
			"msg":    " interal error ",
			"data":   iris.Map{},
		})
	})
}
