package httpmiddleware

import (
	"E7Speed/utils"
	"fmt"

	"github.com/kataras/iris/v12"
)

func HandleRequest(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	fmt.Println(ctx.Path())
	utils.Info(ctx.Path())
	ctx.Next()
}
