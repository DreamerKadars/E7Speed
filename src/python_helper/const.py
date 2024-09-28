from PIL import Image
import numpy as np
import math

SceneImagePath = "./static/scene_image/"

SceneBeibaoName = "beibao"
SceneBeibaoRegion = (140, 78, 220, 110)

SceneIntensifyName = "intensify"
SceneIntensifyRegion = (60, 20, 170, 52)

SceneList = [SceneBeibaoName, SceneIntensifyName]
SceneRegionList = [SceneBeibaoRegion, SceneIntensifyRegion]

def GetSceneFeatureImgLoc(sceneName:str):
    return SceneImagePath+"scene_"+sceneName+"_feature.png"


def mse(imageA, imageB):
    err = np.sum((imageA.astype("float") - imageB.astype("float")) ** 2)
    err /= float(imageA.shape[0] * imageA.shape[1])
    return err


def psnr(imageA, imageB):
    mse_value = mse(imageA, imageB)
    if mse_value == 0:
        return 100
    max_pixel = 255.0
    psnr_value = 20 * math.log10(max_pixel / math.sqrt(mse_value))
    return psnr_value

def IsSameFeature(img: Image, featureLoc: str, region: any):
    cropped_img = img.crop(region)
    # 假设这里有一个已知的参考图片
    reference_img = Image.open(featureLoc)

    # 将图像转换为数组以便进行相似度计算
    cropped_array = np.array(cropped_img)
    reference_array = np.array(reference_img)
    # 计算 MSE 和 PSNR
    mse_value = mse(cropped_array, reference_array)
    psnr_value = psnr(cropped_array, reference_array)

    # 可以根据 PSNR 值设置一个阈值来判断相似度
    threshold = 40
    print("结果")
    print(psnr_value)
    return psnr_value >= threshold
