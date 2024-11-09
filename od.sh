rm -rf /tmp/inputs /tmp/outputs /tmp/*.log 2>/dev/null
mkdir /tmp/inputs 2>/dev/null
export FFREPORT=file=/tmp/ffmpeg.log:level=32
ffmpeg -i $1 -vf fps=$2 /tmp/inputs/image_%4d.png -report
python3 /home/streamer/bins/yolo.py --input_folder=/tmp/inputs/ --output_folder=/tmp/outputs/ > /tmp/yolo.log
ffmpeg -hide_banner -loglevel error -r $2 -i /tmp/outputs/image_%4d.jpg -vcodec libx264 -preset veryfast -crf 25 -movflags +frag_keyframe+empty_moov -f mp4 -y $3 -report