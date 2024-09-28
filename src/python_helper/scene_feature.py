from PIL import Image
from typing import List
from const import *
import json
from json import JSONEncoder,dumps
import os
import cv2
from ultralytics import YOLO
from sanic import Sanic
from sanic.response import json


def crop_and_save(image_path, crop_region, output_path):
    img = Image.open(image_path)
    region = crop_region
    cropped_img = img.crop(region)
    cropped_img.save(output_path)


def GenerateSceneFeature():  # 生成每个界面的场景特征
    for index, scene in enumerate(SceneList):
        print(scene)
        image_path = SceneImagePath+"scene_"+scene+".png"
        output_path = SceneImagePath+"scene_"+scene+"_feature"+".png"
        crop_region = SceneRegionList[index]
        crop_and_save(image_path, crop_region, output_path)


class DetectionResult:
    class_name: str
    x1: int
    y1: int
    x2: int
    y2: int

    def __init__(self, class_name: str = "", x1: int = 0, y1: int = 0, x2: int = 0, y2: int = 0):
        self.class_name = class_name
        self.x1 = x1
        self.y1 = y1
        self.x2 = x2
        self.y2 = y2


class DetectionResultEncoder(JSONEncoder):
    def default(self, obj):
        if isinstance(obj, DetectionResult):
            return {
                "class": obj.class_name,
                "x1": obj.x1,
                "y1": obj.y1,
                "x2": obj.x2,
                "y2": obj.y2
            }
        return super().default(obj)


class SceneModel:
    imgPath: str
    name: str

    def __init__(self, name, imgPath):
        self.imgPath = imgPath
        self.name = name

    @staticmethod
    def IsNowScene(img: Image):
        return False

    @staticmethod
    def GetName():
        return ""

    @staticmethod
    def GetRegion():
        return (0, 0, 0, 0)

    def AnalyseObject(self) -> List[DetectionResult]:
        print("空场景")
        return []


class SceneModelBeibao(SceneModel):
    @staticmethod
    def IsNowScene(img: Image):
        return IsSameFeature(img, GetSceneFeatureImgLoc(SceneBeibaoName), SceneBeibaoRegion)

    @staticmethod
    def GetName():
        return SceneBeibaoName

    @staticmethod
    def GetRegion():
        return SceneBeibaoRegion

    def AnalyseObject(self) -> List[DetectionResult]:
        res: List[DetectionResult] = []
        res.append(DetectionResult(
            class_name="beibao", x1=0, y1=0, x2=0, y2=0))
        try:
            image_path = super().imgPath
            print("开始预测")
            results = model.predict(
                image_path, save=False, imgsz=2048, conf=0.5, iou=0.2)
            print("预测完成")

            result = results[0]
            image_result = cv2.imread(image_path)
            # 读取检测Box
            for i in range(len(result)):
                x1, y1, x2, y2 = result[i].boxes.xyxy[0, 0].item(), result[i].boxes.xyxy[0, 1].item(
                ), result[i].boxes.xyxy[0, 2].item(), result[i].boxes.xyxy[0, 3].item()
                res.append(DetectionResult(
                    class_name=result[i].names[int(
                        result[i].boxes.cls[0])], x1=x1, y1=y1, x2=x2, y2=y2))
                cv2.rectangle(image_result, (int(x1), int(y1)), (int(
                    x2), int(y2)), get_color(int(result[i].boxes.cls[0])), 2)

            cv2.imwrite(os.path.join(image_result_dir,
                        os.path.basename(image_path)), image_result)

            img = cv2.imread(image_path)
            gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
            ret, mask_all_192 = cv2.threshold(
                gray, 192, 255, cv2.THRESH_BINARY)
            cv2.imwrite(os.path.join(image_process_dir_192,
                        os.path.basename(image_path)), mask_all_192)

            ret, mask_all_127 = cv2.threshold(gray, 50, 255, cv2.THRESH_BINARY)
            cv2.imwrite(os.path.join(image_process_dir_127,
                        os.path.basename(image_path)), mask_all_127)
        except Exception as e:
            # 处理异常
            print("发生异常:", e)
        else:
            print("没有发生异常，完成了一次预测")
            print("背包场景")
        return res


class SceneModelIntensify(SceneModel):
    @staticmethod
    def IsNowScene(img: Image):
        return IsSameFeature(img, GetSceneFeatureImgLoc(SceneIntensifyName), SceneIntensifyRegion)

    @staticmethod
    def GetName():
        return SceneIntensifyName

    @staticmethod
    def GetRegion():
        return SceneIntensifyRegion

    def AnalyseObject(self) -> List[DetectionResult]:
        print("强化场景")
        res: List[DetectionResult] = []
        res.append(DetectionResult(
            class_name="intensify", x1=0, y1=0, x2=0, y2=0))
        image_result = cv2.imread(self.imgPath)


        res.append(DetectionResult(
            class_name="mainValueWord", x1=220, y1=218, x2=362, y2=248))
        cv2.rectangle(image_result, (220, 218), (362, 248), 2, 2)

        x1, x2, y1, y2 = 38, 362, 276, 296
        interval = 25.5
        for i in range(4):
            res.append(DetectionResult(
                class_name="subValueWord"+str(i), x1=x1, y1=y1+i*interval, x2=x2, y2=y2+i*interval))
            cv2.rectangle(image_result, (int(x1), int(y1+i*interval)), (int(
                x2), int(y2+i*interval)), 3, 2)
        return res


sceneModelBeibao = SceneModelBeibao("背包", "")
sceneModelIntensify = SceneModelIntensify("强化", "")
SceneModelList: List[SceneModel] = [sceneModelBeibao, sceneModelIntensify]


def GetNowScene(imgPath: str) -> SceneModel:
    res: SceneModel = SceneModel("", "")
    img = open_image(imgPath)
    for index, scene in enumerate(SceneModelList):
        if scene.IsNowScene(img):
            res = scene.__class__(scene.GetName(), imgPath)
            return res
    return res


def open_image(path: str) -> Image.Image:
    try:
        return Image.open(path)
    except Exception as e:
        raise ValueError(f"无法打开图像文件：{path}。错误信息：{e}")


def AnalyseScene(imgPath: str):
    model = GetNowScene(imgPath)
    print(imgPath)
    print(model)
    print(model.GetName())
    return model.AnalyseObject()


def main():
    print(dumps(AnalyseScene(
        "c:\\Users\\RANRAN\\E7Speed\\src\\backend\\screen.png"), cls=DetectionResultEncoder))


if __name__ == "__main__":
    main()
