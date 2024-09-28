import json
import os
import cv2
import pyautogui
from ultralytics import YOLO
from sanic import Sanic
from sanic.response import json

from scene_feature import AnalyseScene

def get_color(number):
    color_dict = {
        0: (255, 0, 0),     # Red
        1: (0, 0, 255),     # Blue
        2: (0, 255, 0),     # Green
        3: (255, 255, 0),   # Yellow
        4: (255, 165, 0),   # Orange
        5: (128, 0, 128),   # Purple
        6: (255, 192, 203), # Pink
        7: (0, 255, 255),   # Cyan
        8: (165, 42, 42),   # Brown
        9: (0, 128, 0),     # Dark Green
        10: (255, 0, 255),  # Magenta
        11: (0, 255, 255),  # Aqua
        12: (128, 128, 128),# Gray
        13: (128, 0, 0),    # Maroon
        14: (0, 128, 128),  # Teal
        15: (128, 128, 0),  # Olive
        16: (192, 192, 192),# Silver
        17: (0, 0, 128),    # Navy
        18: (128, 0, 128),  # Fuchsia
        19: (0, 128, 0),    # Lime
        20: (128, 128, 0),  # Olive
        21: (0, 0, 128),    # Navy
        22: (128, 0, 128),  # Fuchsia
        23: (0, 128, 0),    # Lime
        24: (128, 128, 0)   # Olive
    }

    return color_dict.get(number, (0, 0, 0))  # 默认返回黑色

# 定义图片目录和JSON保存目录
image_dir = "./predict/image"
image_process_dir_127 = "./predict/image_process_127"
image_process_dir_192 = "./predict/image_process_192"
image_result_dir = "./predict/image_result"
json_dir = "./predict/json"

# 创建JSON保存目录

if not os.path.exists(image_dir):
    os.makedirs(image_dir)

if not os.path.exists(image_process_dir_127):
    os.makedirs(image_process_dir_127)

if not os.path.exists(image_process_dir_192):
    os.makedirs(image_process_dir_192)

if not os.path.exists(image_result_dir):
    os.makedirs(image_result_dir)

model = YOLO('./best.pt')

app = Sanic("MyApp")

@app.route("/AnalyseImage")
async def test(request):
    fileLoc = request.args.get('file_loc')
    return json(AnalyseScene(fileLoc))

    beiBaoFlag = False
    intensifyFlag = False
    try:
        beiBaoLoc = pyautogui.locate("./beibao.png",fileLoc)
    except Exception as e:
        print("没有检测到打开背包")
    else:
        beiBaoFlag=True
        print("检测到打开背包")

    try:
        intensifyLoc = pyautogui.locate("./intensify.png", fileLoc)
    except Exception as e:
        print("没有检测到打开强化界面")
    else:
        intensifyFlag = True
        print("检测到打开强化界面")


    if beiBaoFlag:
        resultsJson.append(
            {"class": "beibao", "x1": 0, "y1": 0, "x2": 0, "y2": 0})
        return json(resultsJson)
    elif intensifyFlag:


        #cv2.imwrite(os.path.join(image_result_dir,os.path.basename(image_path)), image_result)
        return json(resultsJson)
    else:
        return json([])


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8000)
