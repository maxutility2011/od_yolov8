rm -rf /tmp/inputs /tmp/outputs 2>/dev/null
mkdir /tmp/inputs 2>/dev/null
ffmpeg -i $1 -vf fps=$2 /tmp/inputs/image_%4d.png -report
python3 /home/streamer/bins/yolo.py --input_folder=/tmp/inputs/ --output_folder=/tmp/outputs/
ffmpeg -hide_banner -loglevel error -r $2 -i /tmp/outputs/image_%4d.jpg -vcodec libx264 -preset veryfast -crf 25 -y /tmp/output.mp4 -report