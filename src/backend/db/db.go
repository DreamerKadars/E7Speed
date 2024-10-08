package db

import (
	"E7Speed/db/operator"
)

// db包负责提供数据库的操作，提供数据访问层和领域层的能力，分别体现在operator和type中,领域层提供的方法会注入进数据访问层，使表层的功能在底层生效

var Operator ChiefOperator

type StaticDataConf struct {
	Dir                        string `yaml:"dir"`
	HeroDataFribbelsFile       string `yaml:"hero_data_fribbels_file"`
	HeroDataFile               string `yaml:"hero_data_file"`
	HeroExtraPanelInfoDataFile string `yaml:"hero_extra_panel_info_data_file"`
	EETypeDataFile             string `yaml:"eeType_data_file"`
}

var staticDataConf StaticDataConf

func SetStaticDataConf(s StaticDataConf) {
	staticDataConf = s
}

func InitDB() {
	err := operator.InitStaticHero(staticDataConf.Dir+staticDataConf.HeroDataFile,
		staticDataConf.Dir+staticDataConf.HeroDataFribbelsFile,
		staticDataConf.Dir+staticDataConf.HeroExtraPanelInfoDataFile,
		staticDataConf.Dir+staticDataConf.EETypeDataFile,
	)
	if err != nil {
		panic(err)
	}
	InitChineseImgHash()
	// client.InitDBConnect()
	Operator = ChiefOperator{}
}

type ChiefOperator struct {
	operator.HeroStaticOperator
}
