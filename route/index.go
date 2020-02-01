package route

import (
	"github.com/kataras/iris/v12"
)

type routeFunc func(*iris.Application)

var routerList []routeFunc

// AddRouter 添加路由
func AddRouter(router routeFunc) {
	routerList = append(routerList, router)
}

// Route 启动路由入口
func Route(app *iris.Application) {
	for _, router := range routerList {
		router(app)
	}
}
