import torch
from ultralytics import YOLO
import cv2
import matplotlib.pyplot as plt
import argparse
import os
from pathlib import Path

print("Current working directory: ", os.getcwd())
os.chdir("/home/streamer/bins/")
print("New working directory: ", os.getcwd())

model = YOLO('yolov8n.pt')
print("yolov8n.pt downloaded")
device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
model = model.to(device)
print("reading arguments")

parser = argparse.ArgumentParser(description='Small object detection with Yolov8.')
parser.add_argument('--loglevel', type=str, default='INFO', help='Set the log level (e.g., DEBUG, INFO, WARNING, ERROR, CRITICAL)')
#parser.add_argument('--trt_engine', type=str, required=True, default='./real_esrgan.engine', help='Path to the vsr inference engine file')
parser.add_argument('--input_folder', type=str, required=True, help='The top level input folder that contains all the input images')
parser.add_argument('--output_folder', type=str, required=True, help='The output folder')

args = parser.parse_args()
print("arguments parsed")
'''
def read_images_from_folder(folder_path):
    images = []
    print("read_images_from_folder")
    for filename in os.listdir(folder_path):
        file_path = os.path.join(folder_path, filename)
        img = cv2.imread(file_path)

        if img is not None:
            images.append(img)
        else:
            print("Failed to load image: ", file_path)
    
    return images
'''

def read_images_from_folder(input_folder):
    images = []
    folder_path = Path(input_folder)
    files = sorted(folder_path.iterdir())

    for file in files:
        image_file = os.path.join(input_folder, file)
        if file.is_file():
            img = cv2.imread(image_file)
            if img is not None:
                images.append(img)
            else:
                print("Failed to load image: ", file.name)

    return images

input_images = read_images_from_folder(args.input_folder)
if len(input_images) == 0:
    print("Failed to read images from: ", dirpath)
    exit(0)

output_folder = args.output_folder
os.makedirs(output_folder, exist_ok=True)

print("Input images loaded. Starting Yolo inference...")
for idx, img in enumerate(input_images):
    #input_img = cv2.resize(img, (320, 192))
    results = model.predict(img, imgsz=320)
    annotated_img = results[0].plot()

    # Display the result
    #plt.imshow(cv2.cvtColor(annotated_img, cv2.COLOR_BGR2RGB))
    #plt.axis('off')
    #plt.show()

    url = output_folder + "/image_" + ("%06d" % (idx+1)) + ".jpg"
    cv2.imwrite(url, annotated_img)
    print("Wrote annotated image: ", url)