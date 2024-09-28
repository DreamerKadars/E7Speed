package model

import "E7Speed/utils"

var model actionModel

const (
	SceneNameIntensify = "Intensify" // 装备强化界面

	ActionOnlyDetect    = "OnlyDetect"    // 装备自动强化
	ActionAutoIntensify = "AutoIntensify" // 装备自动强化
)

func GetModel() actionModel {
	return model
}

type actionModel struct {
	Action                 string            // 当前要执行的行为
	TimeUnit               int               // 时间计数器
	ActionConfig           map[string]string // 当前要执行的行为所配置的选项
	Ctx                    map[string]string // 运行时的环境变量
	SceneName              string            // 界面名称
	MessageErr             string
	DataSceneNameIntensify DataSceneNameIntensify // 装备强化界面名称
}

type DataSceneNameIntensify struct {
}

func (a actionModel) Run() {
	utils.Info("运行model一次：%+v", a)
	// 根据Action，选择对应的策略
	switch a.Action {
	case ActionOnlyDetect:

	}
}
