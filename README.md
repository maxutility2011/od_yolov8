# od_yolov8
Video object detection with Yolo8

- Install python3: 
'''
sudo apt install python3
'''
- Install pip3 and python-venv:
'''
sudo apt update
sudo apt install python3-pip
apt install python3-venv
'''
- Create a Python virtual environment "streamer"
'''
python3 -m venv streamer
''' 
- Activate venv "streamer"
'''
source streamer/bin/activate
'''
When done, deactivate venv "streamer", run
'''
deactivate
'''
- Install torch and dependencies
'''
pip install torch torchvision torchaudio
'''
- Install ultralytics (Yolo)
'''
pip install ultralytics
'''
- Download the source code
'''
git clone https://github.com/maxutility2011/od_yolov8.git
'''
At least 100 GB of disk space is required to hold the packages, model files and intermediate files needed by object detection.

To run object detection on an input, run the following command,
'''
./od.sh [path_to_input] [output_frame_rate]
'''
Please provide the path to the input video and the output video frame rate, e.g., "./od.sh ./input.mp4 25". The shell script will break the input video into a sequence of images with the given frame rate, run Yolov8 inference on the images and generate output images with the detected objects. Finally, it runs ffmpeg to combine the images into an output video. 