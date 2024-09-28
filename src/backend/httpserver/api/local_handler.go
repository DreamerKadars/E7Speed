package api

import (
	"E7Speed/adb"
	"E7Speed/service"
	"E7Speed/utils"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"

	"github.com/kataras/iris/v12"
)

// 用于间接调用本地的一些功能，进行操作和识别

func ParseNowPage(ctx iris.Context) {
	err := adb.ScreenCap()
	if err != nil {
		utils.FailResponse(ctx, err)
		return
	}

	file, err := os.Open("./" + adb.ScreenFileName)
	if err != nil {
		utils.FailResponse(ctx, err)
		return
	}
	defer file.Close()

	tmp := make([]byte, 1e7)
	n, err := file.Read(tmp)
	if err != nil && err != io.EOF {
		utils.FailResponse(ctx, fmt.Errorf("Open file failed,err:%v", err))
		return
	}
	if n >= 1e7 {
		utils.FailResponse(ctx, fmt.Errorf("分辨率过高、无法读取全部图片数据。"))
		return
	}

	file2, err := os.Open("./" + adb.ScreenFileName)
	if err != nil {
		utils.FailResponse(ctx, err)
		return
	}
	defer file2.Close()
	im, _, err := image.DecodeConfig(file2)
	if err != nil {
		utils.FailResponse(ctx, fmt.Errorf("解析图片失败：%+v。", err))
		return
	}
	if im.Width != 1280 || im.Height != 720 {
		utils.FailResponse(ctx, fmt.Errorf("图片分辨率为 %+v*%+v,当前只支持1280*720", im.Width, im.Height))
		return
	}

	ctx.StatusCode(http.StatusOK)

	dir, err := os.Getwd()
	if err != nil {
		utils.FailResponse(ctx, fmt.Errorf("读取目录失败,err:%v", err))
		return
	}
	fileLoc := dir + "\\" + adb.ScreenFileName
	// 发到python上，进行解析
	pythonResp, err := http.Get("http://localhost:8000/AnalyseImage?file_loc=" + fileLoc)
	if err != nil {
		utils.FailResponse(ctx, fmt.Errorf("解析图片失败,err:%v", err))
		return
	}
	var YoloObjects = make([]*service.Object, 0)
	yoloResp := make([]byte, 100000)
	nPython, err := pythonResp.Body.Read(yoloResp)
	if err != nil && err.Error() != "EOF" {
		utils.FailResponse(ctx, fmt.Errorf("解析图片结果失败,err:%v", err))
		return
	}
	yoloResp = yoloResp[0:nPython]

	err = json.Unmarshal(yoloResp, &YoloObjects)
	if err != nil {
		utils.FailResponse(ctx, fmt.Errorf("结构化图片结果失败,err:%v", err))
		return
	}
	utils.Info("%+v", YoloObjects)
	res, err := service.ParseImage(YoloObjects, "screen", ".png", "1")
	if err != nil {
		utils.FailResponse(ctx, fmt.Errorf("解析图片错误,err:%v", err))
		return
	}
	ctx.JSON(res)
}

func AdbConnect(ctx iris.Context) {
	port, err := ctx.URLParamInt("Port")
	if err != nil {
		utils.FailResponse(ctx, err)
		return
	}
	err = adb.AdbConnect(port)
	if err != nil {
		utils.FailResponse(ctx, err)
		return
	}
	err = adb.AdbRoot()
	if err != nil {
		utils.FailResponse(ctx, err)
		return
	}
}

func GetLocalNowImage(ctx iris.Context) {
	adb.ScreenCap()
	file, err := os.Open("./" + adb.ScreenFileName)
	if err != nil {
		utils.FailResponse(ctx, err)
		return
	}
	defer file.Close()

	tmp := make([]byte, 1e7)
	n, err := file.Read(tmp)
	if err != nil && err != io.EOF {
		utils.FailResponse(ctx, fmt.Errorf("Open file failed,err:%v", err))
		return
	}

	ctx.StatusCode(http.StatusOK)
	ctx.Header("Content-Disposition", "attachment; filename="+adb.ScreenFileName)
	ctx.Header("Content-Type", "image/png")
	ctx.Header("Accept-Length", fmt.Sprintf("%d", n))
	ctx.Write(tmp[0:n])
}
