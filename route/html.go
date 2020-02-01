package route

import (
	"github.com/kataras/iris/v12"
)

func init() {
	AddRouter(func(app *iris.Application) {
		app.HandleDir("/", "./")
	})
}
