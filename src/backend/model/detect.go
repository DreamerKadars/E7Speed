package model

import (
	"E7Speed/adb"
	"fmt"
)

type SceneDetect struct {
	ImagePath string
}

func CreateSceneDetect() (*SceneDetect, error) {
	err := adb.ScreenCap()
	if err != nil {
		return nil, err
	}
	return &SceneDetect{
		ImagePath: adb.GetNowImageLoc(),
	}, nil
}

func (s SceneDetect) DetectSceneName() (string, error) {
	return "", fmt.Errorf("未识别的界面")
}

func (s SceneDetect) Detect() {

}
