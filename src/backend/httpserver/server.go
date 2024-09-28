package httpserver

import (
	"E7Speed/httpmiddleware"
	"E7Speed/httpserver/api"
	"E7Speed/httpserver/helper"
	"E7Speed/utils"
	"fmt"
	"os"

	"github.com/kataras/iris/v12"
	"gopkg.in/yaml.v2"
)

var App *iris.Application

type ServiceConf struct {
	JWTConf JWTConf `yaml:"jwt_conf"`
}
type JWTConf struct {
	Secret string `yaml:"secret"`
}

var serviceConf ServiceConf

// 读取服务器鉴权配置
func InitHttpServer() {
	// 读取配置文件中的jwt token
	data, err := os.ReadFile(fmt.Sprintf(utils.Setting.ConfDir + "service_conf.yml"))
	if err != nil {
		panic("load service conf file fail!" + err.Error())
	}
	err = yaml.Unmarshal(data, &serviceConf)
	if err != nil {
		panic("yaml unmarshal fail!" + err.Error())
	}
	httpmiddleware.SetJwtSecret(serviceConf.JWTConf.Secret)
}

func StartHttpServer() {
	utils.InitConfDir()
	// InitHttpServer()
	App = iris.New()

	App.UseRouter(httpmiddleware.HandleRequest)
	apiV1 := App.Party("/api/v1")

	local := apiV1.Party("/local")
	{
		local.Get("/ParseNowPage", api.ParseNowPage)
		local.Get("/AdbConnect", api.AdbConnect)
		local.Get("/GetLocalNowImage", api.GetLocalNowImage)
		local.Get("/hero_template/list", helper.HttpRedirect())
	}
	apiV1.Get("/e7/hero_template/list", helper.HttpRedirect())

	App.Listen(":7766")
}
