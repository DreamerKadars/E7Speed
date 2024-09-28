package utils

const (
	ImageAfterProcess192Path  = "../python_helper/predict/image_process_192/"
	ImageAfterProcess127Path  = "../python_helper/predict/image_process_127/"
	ImageResultPath           = "../python_helper/predict/image_result/"
	ImagePath                 = "../python_helper/predict/image/"
	SaveImagePath             = "../python_helper/predict/save/"
	JsonPath                  = "../python_helper/predict/json/"
	JsonSuffix                = ".json"
	DataFolder                = "../Data/"
	DataHeroStatisticName     = "hero_data.json"
	DataArtifactStatisticName = "artifact_data.json"
	DataEEFribbelsName        = "ee_data_fribbels.json"
	DataHeroImageFolder       = "../Data/HeroImage/"
	DataArtifactImageFolder   = "../Data/ArtifactImage/"

	ClassAtk    = "atk"    // 攻击力
	ClassDefend = "defend" // 防御力
	ClassHp     = "hp"     // 血量
	ClassSpeed  = "speed"  // 速度
	ClassCC     = "cc"     // 暴击率
	ClassCD     = "cd"     // 暴击伤害
	ClassRr     = "rr"     // 效果抵抗
	ClassHr     = "hr"     // 效果命中

	SetAcc       = "acc"       // 命中
	SetAtt       = "att"       // 攻击
	SetCoop      = "coop"      // 夹击
	SetCounter   = "counter"   // 反击
	SetCri       = "cri"       // 暴击
	SetCriDmg    = "cri_dmg"   // 暴伤
	SetImmune    = "immune"    // 免疫
	SetMaxHp     = "max_hp"    // 生命
	SetPenetrate = "penetrate" // 穿透
	SetRage      = "rage"      // 愤怒
	SetRes       = "res"       // 抗性
	SetRevenge   = "revenge"   // 复仇
	SetScar      = "scar"      // 伤口
	SetShield    = "shield"    // 护盾
	SetSpeed     = "speed"     // 速度
	SetTorrent   = "torrent"   // 激流
	SetVampire   = "vampire"   // 吸血
	SetDef       = "def"       // 防御

	EquipLocWeapon   = "weapon"
	EquipLocHelmet   = "helmet"
	EquipLocCuirass  = "cuirass"
	EquipLocNecklace = "necklace"
	EquipLocRing     = "ring"
	EquipLocShoes    = "shoes"
)

var ClassTypeList = []string{
	ClassAtk,
	ClassDefend,
	ClassHp,
	ClassSpeed,
	ClassCC,
	ClassCD,
	ClassRr,
	ClassHr,
}
