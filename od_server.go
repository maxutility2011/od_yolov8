package main

import (
	"flag"
	"fmt"
	//"time"
	//"errors"
	"io/ioutil"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"io"
	"os/exec"
	"od_yolo/job"
)

type Gpu struct {
	Id int
}

type ServerConfig struct {
	Hostname string
	Port string
	Loglevel string
	Gpus []Gpu
}

type AppLogLevel struct { 
	Loglevel string
}

type DetectionJob struct {
	Params job.DetectionParams
	Input_path string
	Output_path string
	State string
	Proc_id string
}

type Detector struct {
	Jobs []DetectionJob
	Gpus []Gpu
}

func readConfig(config_file_path string) ServerConfig {
	configFile, err := os.Open(config_file_path)
	if err != nil {
		slog.Error(err.Error())
	}

	defer configFile.Close()
	config_bytes, _ := ioutil.ReadAll(configFile)
	var config ServerConfig
	json.Unmarshal(config_bytes, &config)
	
	if config.Hostname == "" {
		config.Hostname = "0.0.0.0"
	}

	if config.Port == "" {
		config.Port = "5001"
	}

	if config.Loglevel == "" {
		config.Loglevel = "error"
	}

	return config
}

func contains_string(a []string, v string) bool {
	r := false
	for _, e := range a {
		if v == e {
			r = true
		}
	}

	return r
}

func testPrintLog() {
	slog.Debug("debug")
	slog.Info("info")
	slog.Warn("warn")
	slog.Error("error")
}

var valid_log_levels = []string{"debug", "info", "warn", "error"}
var app_log_level slog.LevelVar
const MaxFileSizeMb = 500

func loglevel_handler(w http.ResponseWriter, r *http.Request) {
	slog.Info("----------------------------------------")
	slog.Info("Received new request to change log level:")
	slog.Info(r.Method, r.URL.Path)

	if !(r.Method == "GET" || r.Method == "PUT") {
		err := "Method = " + r.Method + " is not allowed to " + r.URL.Path
		slog.Error(err)
		http.Error(w, "405 method not allowed\n  Error: "+err, http.StatusMethodNotAllowed)
		return
	}

	if r.Method == "GET" {
		var resp AppLogLevel
		if app_log_level.Level() == slog.LevelDebug {
			resp.Loglevel = "debug"
		} else if app_log_level.Level() == slog.LevelInfo {
			resp.Loglevel = "info"
		} else if app_log_level.Level() == slog.LevelWarn {
			resp.Loglevel = "warn"
		} else if app_log_level.Level() == slog.LevelError {
			resp.Loglevel = "error"
		}

		FileContentType := "application/json"
		w.Header().Set("Content-Type", FileContentType)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	} else if r.Method == "PUT" {
		var l AppLogLevel
		err := json.NewDecoder(r.Body).Decode(&l)
		if err != nil {
			res := "Failed to decode create change_log_level request"
			slog.Error(res)
			http.Error(w, "400 bad request\n  Error: "+res, http.StatusBadRequest)
			return
		}

		if !contains_string(valid_log_levels, l.Loglevel) {
			slog.Error("Invalid log level")
			res := "Invalid log level " + l.Loglevel
			slog.Error(res)
			http.Error(w, "400 bad request\n  Error: "+res, http.StatusBadRequest)
			return
		} else {
			app_log_level.Set(slog.LevelWarn)
			w.WriteHeader(http.StatusOK)
		}
	}
}

func getDetectorArgs(input_path string, output_path string, params job.DetectionParams) []string {
	var args []string
	return args
}

func run_detection(input_path string, params job.DetectionParams) (error, string) {
	output_path := input_path + ".out"

	detectionArgs := getDetectorArgs(input_path, output_path, params)
	detectorCmd := exec.Command("python3", detectionArgs...)
	var err error
	err = nil
	var out []byte
	go func() {
		out, err = detectorCmd.CombinedOutput() // This line blocks when detectorCmd launch succeeds
		if err != nil {
			slog.Error("Errors starting object detector: ", err, ". Detector output: ", string(out))
		}
	}()

	// Wait 100ms before detector fully starts
	/*time.Sleep(100 * time.Millisecond)
	if err != nil {
		slog.Error("Errors starting object detector: ", err, ". Detector output: ", string(out))
		return errors.New("Failed_to_start_detector"), output_path
	}*/

	return nil, output_path
}

func detect_handler(w http.ResponseWriter, r *http.Request) {
	slog.Info("----------------------------------------")
	slog.Info("Received new request to run detection:")
	slog.Info(r.Method, r.URL.Path)

	if !(r.Method == "POST") {
		err := "Method = " + r.Method + " is not allowed to " + r.URL.Path
		slog.Error(err)
		http.Error(w, "405 method not allowed\n  Error: "+err, http.StatusMethodNotAllowed)
		return
	}

	if r.Method == "POST" {
		// curl -F "file=@./output.mp4" -F "params={\"frame_rate\":25,\"reencode_codec\":\"h264\"};type=application/json" http://localhost:5001/detect
		r.Body = http.MaxBytesReader(w, r.Body, MaxFileSizeMb << 20) // Uploaded video file size limit: 500 MB
		err := r.ParseMultipartForm(MaxFileSizeMb << 20) // 500 MB limit for file parsing
		if err != nil {
			fmt.Println(err)
			http.Error(w, "File too large.", http.StatusBadRequest)
			slog.Error("File too large.")
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error retrieving the file.", http.StatusBadRequest)
			slog.Error("Error retrieving the file.")
			return
		}

		defer file.Close()
		dst, err := os.Create(handler.Filename)
		if err != nil {
			http.Error(w, "Unable to save the file.", http.StatusInternalServerError)
			slog.Error("Unable to save the file.")
			return
		}

		defer dst.Close()
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Failed to save the file.", http.StatusInternalServerError)
			slog.Error("Failed to save the file.")
			return
		}

		params_string := r.FormValue("params")
		slog.Info("Params:", params_string)
		var params job.DetectionParams
		//e := json.NewDecoder([]byte(params_string)).Decode(&params)
		e := json.Unmarshal([]byte(params_string), &params)
		if e != nil {
			res := "Failed to decode job params"
			slog.Error("Error happened in JSON marshal. Err: ", e)
			http.Error(w, "400 bad request\n  Error: "+res, http.StatusBadRequest)
			return
		}

		slog.Info("Params:", params)
		run_detection(handler.Filename, params)
	}

	/*
	e1, j := createJob(jspec, warnings)
	if e1 != nil {
		http.Error(w, "500 internal server error\n  Error: ", http.StatusInternalServerError)
		return
	}
		*/



	
}

func (d *Detector) initialize(config ServerConfig) {
	for _, g := range config.Gpus {
		d.Gpus = append(d.Gpus, g)
	}
}

func main() {
	configPtr := flag.String("config", "", "config file path")
	flag.Parse()

	var config_file_path string
	if *configPtr != "" {
		config_file_path = *configPtr
	} else {
		config_file_path = "config.json"
	}

	server_config := readConfig(config_file_path)
	if server_config.Loglevel == "" {
		app_log_level.Set(slog.LevelError)
	} else if server_config.Loglevel == "debug" {
		app_log_level.Set(slog.LevelDebug)
	} else if server_config.Loglevel == "info" {
		app_log_level.Set(slog.LevelInfo)
	} else if server_config.Loglevel == "warn" {
		app_log_level.Set(slog.LevelWarn)
	} else if server_config.Loglevel == "error" {
		app_log_level.Set(slog.LevelError)
	} else {
		fmt.Printf("Unknown log level: %s, use the least verbose level: error. Valid levels are: debug, info, warn and error (ordered in decreasing verbosity).\n", server_config.Loglevel)
		app_log_level.Set(slog.LevelError)
	}

	logfile, err := os.Create("od_server.log")
	if err != nil {
    	panic(err)
	}

	h := slog.NewTextHandler(logfile, &slog.HandlerOptions{Level: &app_log_level})
	slog.SetDefault(slog.New(h))

	var detector Detector 
	detector.initialize(server_config)
	server_addr := server_config.Hostname + ":" + server_config.Port
	
	fmt.Printf("API server listening on: %s\n", server_addr)
	http.HandleFunc("/loglevel", loglevel_handler)
	http.HandleFunc("/detect", detect_handler)
	http.ListenAndServe(server_addr, nil)
}