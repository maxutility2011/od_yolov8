import torch
from ultralytics import YOLO
import cv2
import matplotlib.pyplot as plt
import argparse
import os

model = YOLO('yolov8n.pt')
device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
model = model.to(device).eval()

parser = argparse.ArgumentParser(description='Small object detection with Yolov8.')
parser.add_argument('--loglevel', type=str, default='INFO', help='Set the log level (e.g., DEBUG, INFO, WARNING, ERROR, CRITICAL)')
#parser.add_argument('--trt_engine', type=str, required=True, default='./real_esrgan.engine', help='Path to the vsr inference engine file')
parser.add_argument('--input_folder', type=str, required=True, help='The top level input folder that contains all the input images')
parser.add_argument('--output_folder', type=str, required=True, help='The output folder')

args = parser.parse_args()

def read_images_from_folder(folder_path):
    images = []

    for filename in os.listdir(folder_path):
        file_path = os.path.join(folder_path, filename)
        img = cv2.imread(file_path)

        if img is not None:
            images.append(img)
        else:
            logging.error("Failed to load image: %s", file_path)
    
    return images

input_images = read_images_from_folder(args.input_folder)
if len(input_images) == 0:
    logger.error("Failed to read images from %s", dirpath)
    exit(0)

output_folder = args.output_folder
os.makedirs(output_folder, exist_ok=True)

for idx, img in enumerate(input_images):
    input_img = cv2.resize(img, (1280, 1280))
    results = model(input_img, conf=0.25)
    annotated_img = results[0].plot()

    # Display the result
    #plt.imshow(cv2.cvtColor(annotated_img, cv2.COLOR_BGR2RGB))
    #plt.axis('off')
    #plt.show()

    url = output_folder + "/image_" + ("%04d" % (idx+1)) + ".jpg"
    cv2.imwrite(url, annotated_img)

    # If you want to save the result
    #cv2.imwrite('output_image.jpg', annotated_img)