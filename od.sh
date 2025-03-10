# Parameters
# 1. Input file path
# 2. Frame rate
# 3. Output directory
# 4. Per-segment temporary directory
# 5. Output filename
# 6. SIDX timescale

# Creating input image subdirectory
seg_dir="$3/$4"
mkdir $seg_dir 2>/dev/null 
input_image_dir="$seg_dir/inputs"
mkdir $input_image_dir 2>/dev/null

# Configure ffmpeg image_converter log
image_converter_log="$3/image_converter.log"
export FFREPORT=file=$image_converter_log:level=32

# Convert the input video ($1) to Yolo input images at the given frame rate ($2)
input_images="$input_image_dir/image_%6d.png"
ffmpeg -i $1 -vf fps=$2 $input_images -report

# Run Yolo inference (i.e., object detection on the input images)
output_image_dir="$seg_dir/outputs" # Output directory will be created by Yolo script, no need to mkdir
yolo_log="$seg_dir/yolo.log"
python3 /home/streamer/bins/yolo.py --input_folder=$input_image_dir --output_folder=$output_image_dir > $yolo_log

# Re-encode the Yolo output images (annotated) to the given output video ($4) at the given frame rate ($2)
output_images="$output_image_dir/image_%6d.jpg"
# Configure ffmpeg re-encoder log
#reencoder_log="$3/reencoder.log"
#export FFREPORT=file=$reencoder_log:level=32
# Re-encoder uses libx264, preset=veryfast, crf=25. 
# Re-encoder outputs a fragmented mp4 segment (e.g., seg_1.detected) including the init section (FTYP + MOOV) and one MOOF atom (data section)
# The re-encoder output is not immediately uploadable. Worker_transcoder still needs to split the segment into two upload candidates, e.g.,
# - an init segment
# - an media data segment
# TODO: The live jobs have the configuration for re-encoder (e.g., codec, preset, crf). We should use those values instead.
# Example: ffmpeg -hide_banner -loglevel error -r 5 -i /tmp/output_2108d5f3-a590-4a3e-a3e6-4231c0ee076d/video_150k/seg_1/outputs/image_%6d.jpg -vcodec libx264 -preset veryfast -crf 25 -video_track_timescale 90000 -movflags +frag_keyframe+empty_moov -f mp4 -y 1.mp4
ffmpeg -hide_banner -loglevel error -r $2 -i $output_images -vcodec $7 -preset $8 -crf $9 -video_track_timescale $6 -movflags +frag_keyframe+empty_moov -f mp4 -y $5 -report

# Delete intermediate files
rm -rf $input_image_dir 2>/dev/null 
rm -rf $output_image_dir 2>/dev/null 