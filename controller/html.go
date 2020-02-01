package controller

import (
	"github.com/kataras/iris/v12"
)

// Get "/"
func Get(ctx iris.Context) {
	ctx.HTML("<h1>Hello LiveServer</h1>")
}
