package adb

import (
	"E7Speed/utils"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
)

var (
	AdbExeLoc        = `C:\Program Files\Netease\MuMuPlayer-12.0\shell\adb.exe`
	ScreenFileName   = "screen.png"
	ScreenFilePrefix = "/data/"
)

func GetNowImageLoc() string {
	return ScreenFilePrefix + ScreenFileName
}

// 连接到指定的端口上
func AdbConnect(Port int) error {
	cmd := exec.Command(AdbExeLoc, "connect", "127.0.0.1:"+strconv.Itoa(Port))
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}
	return nil
}

// 获取Root权限
func AdbRoot() error {
	cmd := exec.Command(AdbExeLoc, "root")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}
	return nil
}

// 截图并尝试保存到本地
func ScreenCap() (err error) {
	cmd := exec.Command(AdbExeLoc, "shell", "screencap", GetNowImageLoc())
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		utils.Error(fmt.Sprint(err) + ": " + stderr.String())
		return
	}

	cmd = exec.Command(AdbExeLoc, "pull", GetNowImageLoc(), "./")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		utils.Error(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	return
}
