rm -rf ./inputs ./outputs 2>/dev/null
mkdir ./inputs 2>/dev/null
ffmpeg -i $1 -vf fps=$2 ./inputs/image_%4d.png
python3 yolo.py --input_folder=./inputs/ --output_folder=./outputs/
ffmpeg -hide_banner -loglevel error -r $2 -i ./outputs/image_%4d.jpg -vcodec libx264 -preset veryfast -crf 25 -y ./output.mp4