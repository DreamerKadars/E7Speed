package service

import (
	"E7Speed/db"
	"E7Speed/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	ClassEquipment      = "equipment" // 完整装备
	ClassLevel85        = "85"        // 85级装备
	dpiMin              = 70
	dpiMax              = 2400
	numLengthParam      = 2.10
	numLengthParamLarge = 3.0

	equipLocE1 = "e1"
	equipLocE2 = "e2"
	equipLocE3 = "e3"
	equipLocE4 = "e4"
	equipLocE5 = "e5"
	equipLocE6 = "e6"

	ParseSetAcc       = "hr"        // 命中
	ParseSetAtt       = "atk"       // 攻击
	ParseSetCoop      = "coop"      // 夹击
	ParseSetCounter   = "counter"   // 反击
	ParseSetCri       = "cc"        // 暴击
	ParseSetCriDmg    = "cd"        // 暴伤
	ParseSetImmune    = "immune"    // 免疫
	ParseSetMaxHp     = "hp"        // 生命
	ParseSetPenetrate = "penetrate" // 穿透
	ParseSetRage      = "rage"      // 愤怒
	ParseSetRes       = "rr"        // 抗性
	ParseSetRevenge   = "revenge"   // 复仇
	ParseSetScar      = "scar"      // 伤口
	ParseSetShield    = "shield"    // 护盾
	ParseSetSpeed     = "speed"     // 速度
	ParseSetTorrent   = "torrent"   // 激流
	ParseSetVampire   = "vampire"   // 吸血
	ParseSetDef       = "defend"    // 防御

	EquipColorRed    = "red"
	EquipColorPurple = "purple"

	equipLocRedKey    = "redE"
	equipLocPurpleKey = "purpleE"
)

var ConvertSetType map[string]string = map[string]string{
	ParseSetAcc:       utils.SetAcc,
	ParseSetAtt:       utils.SetAtt,
	ParseSetCoop:      utils.SetCoop,
	ParseSetCounter:   utils.SetCounter,
	ParseSetCri:       utils.SetCri,
	ParseSetCriDmg:    utils.SetCriDmg,
	ParseSetImmune:    utils.SetImmune,
	ParseSetMaxHp:     utils.SetMaxHp,
	ParseSetPenetrate: utils.SetPenetrate,
	ParseSetRage:      utils.SetRage,
	ParseSetRes:       utils.SetRes,
	ParseSetRevenge:   utils.SetRevenge,
	ParseSetScar:      utils.SetScar,
	ParseSetShield:    utils.SetShield,
	ParseSetSpeed:     utils.SetSpeed,
	ParseSetTorrent:   utils.SetTorrent,
	ParseSetVampire:   utils.SetVampire,
	ParseSetDef:       utils.SetDef,
}

var ConvertLocValue map[string]string = map[string]string{
	equipLocE1: utils.EquipLocWeapon,
	equipLocE2: utils.EquipLocHelmet,
	equipLocE3: utils.EquipLocCuirass,
	equipLocE4: utils.EquipLocNecklace,
	equipLocE5: utils.EquipLocRing,
	equipLocE6: utils.EquipLocShoes,
}

type Equipment struct {
	ID string // 随机数ID
	Object
	CC             int       //暴击率
	CD             int       //暴击伤害
	Atk            int       //攻击白
	AtkPercent     int       //攻击百分比
	Speed          int       //速度
	Hp             int       //血量
	HpPercent      int       //血量百分比
	RR             int       //效果抵抗
	Hr             int       //效果命中
	Defend         int       //防御力
	DefendPercent  int       //防御力百分比
	Level          int       // 等级
	UpgradeLevel   int       // 强化等级
	UpgradePercent []float64 // 升级到各个级别装备的概率
	AnchorIndex    int       // 主属性锚点下标
	Objects        []Object  // 包含在内的检测物体
	MainType       string    // 主属性类型
	MainValue      int       // 主属性数值
	Set            string    // 套装类型
	EquipLoc       string    // 装备位置
	EquipColor     string    // 装备颜色，只支持红色和紫色
	FourLocClass   string    // 紫装的第四个副属性是特殊的
}

type Object struct {
	Class   string `json:"class_name"` // 类别
	X1      float32
	Y1      float32
	X2      float32
	Y2      float32
	Value   int
	Percent bool
}

// 因为效果不好，需要进行纠正
func CorrectYoloObjects(YoloObjects []*Object) []*Equipment {
	equipment := make([]*Equipment, 0)
	otherObject := make([]*Object, 0)
	levelObject := make([]*Object, 0)

	for _, object := range YoloObjects {
		if object.Class == ClassEquipment {
			equipment = append(equipment, &Equipment{
				Object: *object,
			})
		} else if object.Class == ClassLevel85 {
			levelObject = append(levelObject, object)
		} else {
			otherObject = append(otherObject, object)
		}
	}

	// 属性归属
	for _, objectOther := range otherObject {
		for _, equip := range equipment {
			if objectOther.X2 <= equip.X2+(equip.X2-equip.X1)/10 &&
				objectOther.Y2 <= equip.Y2+(equip.Y2-equip.Y1)/10 &&
				objectOther.X1 >= equip.X1-(equip.X2-equip.X1)/10 &&
				objectOther.Y1 >= equip.Y1-(equip.Y2-equip.Y1)/10 {
				equip.Objects = append(equip.Objects, *objectOther)
			}
		}
	}

	// 等级判断
	for _, levelObjectTemp := range levelObject {
		for _, equip := range equipment {
			if levelObjectTemp.X2 <= equip.X2+(equip.X2-equip.X1)/10 &&
				levelObjectTemp.Y2 <= equip.Y2+(equip.Y2-equip.Y1)/10 &&
				levelObjectTemp.X1 >= equip.X1-(equip.X2-equip.X1)/10 &&
				levelObjectTemp.Y1 >= equip.Y1-(equip.Y2-equip.Y1)/10 {
				equip.Level = 85 // 当前只有85级装备
			}
		}
	}

	// 删除没有两条属性以上的装备
	newEquip := make([]*Equipment, 0)
	for _, equip := range equipment {
		if len(equip.Objects) > 1 {
			newEquip = append(newEquip, equip)
		}
	}
	equipment = newEquip

	// 寻找每个装备的主属性位置
	for _, equip := range equipment {
		// 装备的主属性被认为是锚点，即坐标y值最小的那个，和坐标系是反的
		var indexAnchor int = -1
		var minY float32 = 100000
		for index, object := range equip.Objects {
			if object.Y2 < minY {
				minY = object.Y2
				indexAnchor = index
			}
		}
		equip.AnchorIndex = indexAnchor
	}

	return equipment
}

func IsClassType(class string) bool {
	switch class {
	case utils.ClassAtk, utils.ClassHp, utils.ClassDefend, utils.ClassSpeed, utils.ClassCC, utils.ClassCD, utils.ClassRr, utils.ClassHr:
		return true
	}
	return false
}

func YoloReadImage(imagePath string) ([]*Object, error) {
	res := make([]*Object, 0)
	// 发给python server，用来进行识别
	client := &http.Client{}
	// 发送 GET 请求
	response, err := client.Get("http://localhost:8000/AnalyseImage?file_loc=" + imagePath)
	if err != nil {
		fmt.Println("GET request failed:", err)
		return nil, err
	}
	defer response.Body.Close()
	// 读取响应内容
	body := make([]byte, 10024)
	_, err = response.Body.Read(body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

const (
	ImageTypeBeiBao    = "beibao"
	ImageTypeIntensify = "intensify"
)

const (
	// 在1280*720的条件下，进行的强化等级位置的计算
	Equip1280720UpgradeLevelLocX1 = 100
	Equip1280720UpgradeLevelLocX2 = 135
	Equip1280720UpgradeLevelLocY1 = 83
	Equip1280720UpgradeLevelLocY2 = 103

	// 在1280*720的条件下，套装信息的位置的x上限和y上限
	Equip1280720SetLocMaxX1 = 210
	Equip1280720SetLocMaxY1 = 208

	// 在1280*720的条件下，主属性信息的位置的x上限和y上下限
	Equip1280720MainTypeLocMax1  = 210
	Equip1280720MainTypeLocMinY1 = 208
	Equip1280720MainTypeLocMaxY1 = 238

	// 在1280*720的条件下，强化时的占位长度
	Equip1280720SubTypeIntensifyLen = 55

	MainValueWordClassName = "mainValueWord"
	SubValueWordClassName  = "subValueWord"
)

func ParseImageIntensify(YoloObjects []*Object, imageFileName, imageSuffix string, ID string) ([]*Equipment, error) {
	indexImage := 1
	equip := Equipment{}
	equip.Level = 88

	// 计算强化等级
	value, _, _ := CalculateNum(indexImage, utils.ImageAfterProcess192Path, imageFileName, imageSuffix, true, Equip1280720UpgradeLevelLocX1, Equip1280720UpgradeLevelLocY1, Equip1280720UpgradeLevelLocX2, Equip1280720UpgradeLevelLocY2)
	equip.UpgradeLevel = value

	indexImage++

	// 因为是正在强化，只会生成一个结果
	for _, object := range YoloObjects {
		if object.X1 < Equip1280720SetLocMaxX1 && object.Y1 < Equip1280720SetLocMaxY1 {
			// 左上角的被判断是装备套装
			if set, exist := ConvertSetType[object.Class]; exist {
				equip.Set = set
			}
			// 左上角的被判断是装备等级
			if object.Class == ClassLevel85 {
				equip.Level = 85
			}
			// 判断装备颜色
			if object.Class == equipLocRedKey {
				equip.EquipColor = EquipColorRed
			} else if object.Class == equipLocPurpleKey {
				equip.EquipColor = EquipColorPurple
			}

			//判断装备位置
			if locStr, exist := ConvertLocValue[object.Class]; exist {
				equip.EquipLoc = locStr
			}
		}

		// 主属性范围内的是主属性类型
		if object.X1 < Equip1280720MainTypeLocMax1 && object.Y1 < Equip1280720MainTypeLocMaxY1 && object.Y1 > Equip1280720MainTypeLocMinY1 {
			if IsClassType(object.Class) {
				equip.MainType = object.Class + equip.MainType // 有可能原本是percent
			}
		}

		if object.Class == MainValueWordClassName {
			// 求主属性值
			indexImage++
			value, percent, err := CalculateNum(indexImage, utils.ImageAfterProcess127Path, imageFileName, imageSuffix, false, int(object.X1), int(object.Y1), int(object.X2), int(object.Y2))
			if err != nil {
				utils.Error("call CalculateNum fail:%+v", err)
			}
			utils.Info("ParseWordAndNum res :%+v 百分比：%+v", value, percent)
			if percent && (object.Class == utils.ClassAtk || object.Class == utils.ClassDefend || object.Class == utils.ClassHp) {
				equip.MainType += "Percent"
			}
			equip.MainValue = value
		} else if strings.HasPrefix(object.Class, SubValueWordClassName) {

			// 求副属性
			indexImage++
			subTypeNow, err := ParseChinese(indexImage, utils.ImageAfterProcess127Path, imageFileName, imageSuffix, int(object.X1), int(object.Y1), int((object.X2+object.X1)/2), int(object.Y2))
			if err != nil {
				utils.Error("call ParseChinese fail:%+v", err)
			} else {
				utils.Info("call ParseChinese :%+v %+v", subTypeNow, err)
			}

			// 求副属性值
			indexImage++
			// 有两种可能,强化过的和没强化过的，分开检测
			value, percent, err1 := CalculateNum(indexImage, utils.ImageAfterProcess127Path, imageFileName, imageSuffix, false, int((object.X2+object.X1)/2), int(object.Y1), int(object.X2-Equip1280720SubTypeIntensifyLen), int(object.Y2))
			var err2 error
			if err1 != nil {
				value, percent, err2 = CalculateNum(indexImage, utils.ImageAfterProcess127Path, imageFileName, imageSuffix, false, int(object.X2-Equip1280720SubTypeIntensifyLen), int(object.Y1), int(object.X2), int(object.Y2))
				if err2 != nil {
					utils.Error("call CalculateNum 两次都未查询到结果 fail:%+v %+v", err1, err2)
				}
				utils.Info("call CalculateNum 第二个位置 :%+v %+v", value, percent)
			} else {
				utils.Info("call CalculateNum 第一个位置:%+v %+v", value, percent)
			}

			if strings.TrimPrefix(object.Class, SubValueWordClassName) == "3" {
				equip.FourLocClass = subTypeNow
			}

			switch subTypeNow {
			case utils.ClassAtk:
				if percent {
					equip.AtkPercent = value
				} else {
					equip.Atk = value
				}
			case utils.ClassHp:
				if percent {
					equip.HpPercent = value
				} else {
					equip.Hp = value
				}
			case utils.ClassDefend:
				if percent {
					equip.DefendPercent = value
				} else {
					equip.Defend = value
				}
			case utils.ClassSpeed:
				equip.Speed = value
			case utils.ClassCC:
				equip.CC = value
			case utils.ClassCD:
				equip.CD = value
			case utils.ClassRr:
				equip.RR = value
			case utils.ClassHr:
				equip.Hr = value
			}
		}
	}
	utils.Info("%+v", equip)
	return []*Equipment{&equip}, nil
}

func ParseImageBeiBao(YoloObjects []*Object, imageFileName, imageSuffix string, ID string) ([]*Equipment, error) {
	var err error
	res := CorrectYoloObjects(YoloObjects)
	indexImage := 1
	for _, equip := range res {
		// 装备的主属性被认为是锚点，即右上角坐标y值最大的那个
		var indexAnchor = equip.AnchorIndex

		leftX := equip.Objects[indexAnchor].X1 - (equip.Objects[indexAnchor].X2-equip.Objects[indexAnchor].X1)/5

		// 计算第四个副属性的位置，x最大，y最大的就是
		FourLocClass, FourLocV := "", float32(-100000)

		// 锚点右侧的属性需要参与计算数值，代表副词条
		for index, object := range equip.Objects {
			if index == indexAnchor {
				// 主属性
				value, percent, err := CalculateNum(indexImage, utils.ImageAfterProcess127Path, imageFileName, imageSuffix, false, int(object.X2), int(object.Y1), int(object.X2+numLengthParamLarge*(object.X2-object.X1)), int(object.Y2))
				if err != nil {
					continue
				}
				equip.MainType = object.Class
				if percent && (object.Class == utils.ClassAtk || object.Class == utils.ClassDefend || object.Class == utils.ClassHp) {
					equip.MainType += "Percent"
				}
				equip.MainValue = value
				continue
			}

			percent := false
			if object.X1 > leftX {
				value := 0
				value, percent, err = CalculateNum(indexImage, utils.ImageAfterProcess127Path, imageFileName, imageSuffix, false, int(object.X2), int(object.Y1), int(object.X2+numLengthParam*(object.X2-object.X1)), int(object.Y2))
				indexImage++
				if err != nil {
					continue
				}
				switch object.Class {
				case utils.ClassAtk:
					if percent {
						equip.AtkPercent = value
					} else {
						equip.Atk = value
					}
				case utils.ClassHp:
					if percent {
						equip.HpPercent = value
					} else {
						equip.Hp = value
					}
				case utils.ClassDefend:
					if percent {
						equip.DefendPercent = value
					} else {
						equip.Defend = value
					}
				case utils.ClassSpeed:
					equip.Speed = value
				case utils.ClassCC:
					equip.CC = value
				case utils.ClassCD:
					equip.CD = value
				case utils.ClassRr:
					equip.RR = value
				case utils.ClassHr:
					equip.Hr = value
				}
			}

			if index != indexAnchor {
				// 非主属性的副属性
				if IsClassType(object.Class) {
					if object.X1+object.Y1 > FourLocV {
						FourLocV = object.X1 + object.Y1
						FourLocClass = object.Class
						if percent {
							FourLocClass += "Percent"
						}
					}
				}
			}

			// 锚点左侧的被判断是装备套装
			if object.X2 < equip.Objects[indexAnchor].X1 {
				if set, exist := ConvertSetType[object.Class]; exist {
					equip.Set = set
				}
			}
			//判断装备位置
			if locStr, exist := ConvertLocValue[object.Class]; exist {
				equip.EquipLoc = locStr
			}
			// 判断装备颜色
			if object.Class == equipLocRedKey {
				equip.EquipColor = EquipColorRed
			} else if object.Class == equipLocPurpleKey {
				equip.EquipColor = EquipColorPurple
			}
		}

		equip.FourLocClass = FourLocClass
		// 读取装备强化等级
		x1, y1, x2, y2 := GetUpgradeLevelLoc(equip.X1, equip.Y1, equip.X2, equip.Y2)
		equip.UpgradeLevel, _, err = CalculateNum(indexImage, utils.ImageAfterProcess192Path, imageFileName, imageSuffix, true, x1, y1, x2, y2)
		indexImage++
		if err != nil {
			utils.Error("%+v", err)
		}
	}
	return res, nil
}

type ParseImageRes struct {
	Equips []*Equipment
	Mode   string
}

// 解析图片文件并且得出一系列信息
func ParseImage(YoloObjects []*Object, imageFileName, imageSuffix string, ID string) (ParseImageRes, error) {
	res := make([]*Equipment, 0)
	var err error

	imageType := ""
	// 暂时支持两种模式，通过两个常量来标识，如果带有对应的常量就是用相应的功能
	for _, o := range YoloObjects {
		if o.Class == ImageTypeBeiBao || o.Class == ImageTypeIntensify {
			imageType = o.Class
			break
		}
	}
	switch imageType {
	case ImageTypeBeiBao:
		res, err = ParseImageBeiBao(YoloObjects, imageFileName, imageSuffix, ID)
	case ImageTypeIntensify:
		res, err = ParseImageIntensify(YoloObjects, imageFileName, imageSuffix, ID)
	}

	return ParseImageRes{
		Equips: res,
		Mode:   imageType,
	}, err
}

// 从中文描述转到类型
func GetSubTypeFromChinese(chineseInfo string) (string, error) {
	if true {
		return utils.ClassAtk, nil
	}
	if true {
		return utils.ClassDefend, nil
	}
	if true {
		return utils.ClassHp, nil
	}
	if true {
		return utils.ClassSpeed, nil
	}
	if true {
		return utils.ClassCD, nil
	}
	if true {
		return utils.ClassRr, nil
	}
	if true {
		return utils.ClassHr, nil
	}

	return "", nil
}

// 计算图片中文的副属性
func ParseChinese(index int, imagePath, imageFileName, imageSuffix string, x1, y1, x2, y2 int) (res string, err error) {
	var targetImage string
	var resResult string
	defer func() {
		utils.Info("OCR： 生成图像文件[%s],ocr识别结果：[%s],是否有错误[%+v]", targetImage, resResult, err)
	}()
	targetImage = imagePath + imageFileName + "_temp_" + strconv.Itoa(index) + imageSuffix
	err = Cut(imagePath+imageFileName+imageSuffix, targetImage, x1, y1, x2, y2)
	if err != nil {
		return resResult, err
	}
	queryFile, err := os.Open(targetImage)
	if err != nil {
		return resResult, err
	}

	imgQuery, err := png.Decode(queryFile)
	if err != nil {
		return resResult, err
	}

	minDistance := 1000
	for class, imgTemp := range db.ChineseSubTypeImageMap {
		distance := db.DistanceImage(imgQuery, *imgTemp)
		if err != nil {
			return resResult, err
		}
		if distance < minDistance {
			minDistance = distance
			resResult = class
		}
		utils.Info("[%+v]class[%+v] between images: %d\n", targetImage, class, distance)
	}

	if minDistance > 400 {
		resResult = ""
	}

	return resResult, err
}

// 计算图片数字+符号
func ParseNumAndSymbol(index int, imagePath, imageFileName, imageSuffix string, x1, y1, x2, y2 int) (res string, err error) {
	var targetImage string
	var ocrResult string
	defer func() {
		utils.Info("OCR： 生成图像文件[%s],ocr识别结果：[%s],是否有错误[%+v]", targetImage, ocrResult, err)
	}()
	targetImage = imagePath + imageFileName + "_temp_" + strconv.Itoa(index) + imageSuffix
	err = Cut(imagePath+imageFileName+imageSuffix, targetImage, x1, y1, x2, y2)
	if err != nil {
		return ocrResult, err
	}

	dpi := int((y2 - y1) * 9 / 10)

	if dpi < dpiMin {
		dpi = dpiMin
	}
	if dpi > dpiMax {
		dpi = dpiMax
	}
	cmd := exec.Command("tesseract", targetImage, "stdout", "-l", "E7", "--dpi", strconv.Itoa(dpi), "-c", "tessedit_char_whitelist=0123456789%+()", "--psm", "8")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	ocrResult = out.String()

	return ocrResult, err
}

// 计算图片中一个副属性的数值,强制需要+号时必须要识别到+号
func CalculateNum(index int, imagePath, imageFileName, imageSuffix string, mustPlus bool, x1, y1, x2, y2 int) (value int, percent bool, err error) {
	var targetImage string
	var ocrResult string
	defer func() {
		utils.Info("OCR： 生成图像文件[%s],ocr识别结果：[%s],转义结果[数字:%d 百分比;%+v],是否有错误[%+v]", targetImage, ocrResult, value, percent, err)
	}()
	targetImage = imagePath + imageFileName + "_temp_" + strconv.Itoa(index) + imageSuffix
	err = Cut(imagePath+imageFileName+imageSuffix, targetImage, x1, y1, x2, y2)
	if err != nil {
		return 0, false, err
	}

	dpi := int((y2 - y1) * 9 / 10)

	if dpi < dpiMin {
		dpi = dpiMin
	}
	if dpi > dpiMax {
		dpi = dpiMax
	}
	cmd := exec.Command("tesseract", targetImage, "stdout", "-l", "E7", "--dpi", strconv.Itoa(dpi), "-c", "tessedit_char_whitelist=0123456789%+", "--psm", "8")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	ocrResult = out.String()
	result := ocrResult
	result = strings.Replace(result, "020", "%", -1)
	if strings.Contains(result, "%") {
		percent = true
	}
	// 如果存在百分号，删除%号后面的字符
	if strings.LastIndex(result, "%") != -1 {
		result = result[:strings.LastIndex(result, "%")]
	}

	result = strings.ReplaceAll(result, "%", "")
	result = strings.ReplaceAll(result, "\n", "")
	result = strings.ReplaceAll(result, "\f", "")
	result = strings.ReplaceAll(result, "\r", "")
	result = strings.ReplaceAll(result, " ", "")

	if len(result) == 0 {
		return 0, false, fmt.Errorf("未解析到数字")
	}
	if mustPlus && !strings.Contains(result, "+") {
		return 0, false, nil
	}
	// 如果存在+号，删除+号前的字符
	if strings.LastIndex(result, "+") != -1 {
		result = result[strings.LastIndex(result, "+"):]
	}
	result = strings.TrimPrefix(result, "+")
	value, err = strconv.Atoi(result)
	return value, percent, err
}

func Cut(source, target string, x1, y1, x2, y2 int) error {
	img, tt, err := loadImage(source)
	if err != nil {
		return err
	}
	if !rectIsInRect(x1, y1, x2, y2, img.Bounds()) {
		return fmt.Errorf("rectangele not in target image")
	}
	// 图片文件解码
	rgbImg := img.(*image.Gray)
	subImg := rgbImg.SubImage(image.Rect(x1, y1, x2, y2)).(*image.Gray) //图片裁剪x0 y0 x1 y1
	return saveImage(target, subImg, 100, tt)
}

func rectIsInRect(x1, y1, x2, y2 int, ret image.Rectangle) bool {
	var p = image.Pt(x1, y1)
	var p1 = image.Pt(x2, y2)
	return pointInRect(p, ret) && pointInRect(p1, ret)

}

// 判断点是否在图片像素的矩形框内
func pointInRect(p image.Point, ret image.Rectangle) bool {
	return ret.Min.X <= p.X && p.X <= ret.Max.X &&
		ret.Min.Y <= p.Y && p.Y <= ret.Max.Y
}

// 加载图片
func loadImage(path string) (image.Image, string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()
	return image.Decode(file)
}

// 保存图片
func saveImage(path string, subImg image.Image, quality int, tt string) error {
	f, err := os.Create(path) //创建文件，会自动覆盖
	if err != nil {
		return err
	}
	defer f.Close() //关闭文件
	var opt jpeg.Options
	opt.Quality = quality

	switch tt {
	case "jpeg":
		return jpeg.Encode(f, subImg, &opt)
	case "png":
		return png.Encode(f, subImg)
	default:
	}
	return nil
}

// 0 0.25968 0.47286 1
// 0 0 0.24633 1
func GetUpgradeLevelLoc(x1, y1, x2, y2 float32) (newX1, newY1, newX2, newY2 int) {
	newX1 = int(x1 + 0.28968*(x2-x1))
	newY1 = int(y1)
	newX2 = int(x1 + 0.42286*(x2-x1))
	newY2 = int(y1 + 0.21633*(y2-y1))
	return
}

// 0 0.05943 0.17700 1
// 0 0.10850 0.33137 1
func GetLevelLoc(x1, y1, x2, y2 float32) (newX1, newY1, newX2, newY2 int) {
	newX1 = int(x1 + 0.05943*(x2-x1))
	newY1 = int(y1 + 0.10850*(y2-y1))
	newX2 = int(x1 + 0.17700*(x2-x1))
	newY2 = int(y1 + 0.33137*(y2-y1))
	return
}
