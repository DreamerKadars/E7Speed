package db

import (
	"E7Speed/utils"
	"image"
	"image/png"
	"os"
)

var ChineseSubTypeImageMap map[string]*image.Image

func InitChineseImgHash() {
	ChineseSubTypeImageMap = make(map[string]*image.Image)
	for _, class := range utils.ClassTypeList {
		file, err := os.Open(staticDataConf.Dir + class + ".png")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		img, err := png.Decode(file)
		if err != nil {
			panic(err)
		}
		ChineseSubTypeImageMap[class] = &img
	}
}

func DistanceImage(img1, img2 image.Image) int {
	if img1.Bounds().Max.X != img2.Bounds().Max.X || img1.Bounds().Max.Y != img2.Bounds().Max.Y {
		return 1e7
	}
	height := img1.Bounds().Max.Y
	width := img1.Bounds().Max.X
	distance := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r1, g1, b1, a1 := img1.At(x, y).RGBA()
			r2, g2, b2, a2 := img2.At(x, y).RGBA()
			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
				distance += 1
			}
		}
	}
	return distance
}
