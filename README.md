# od_yolov8
Video object detection with Yolo8

To run object detection on an input, run the following command,
'''
./od.sh [path_to_input] [output_frame_rate]
'''
Please provide the path to the input video and the output video frame rate, e.g., "./od.sh ./input.mp4 25". The shell script will break the input video into a sequence of images with the given frame rate, run Yolov8 inference on the images and generate output images with the detected objects. Finally, it runs ffmpeg to combine the images into an output video. 