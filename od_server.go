package main

import (
	"flag"
	"fmt"
	//"errors"
	"io/ioutil"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	//"os/exec"
)

type ServerConfig struct {
	Hostname string
	Port string
	Loglevel string
}

type AppLogLevel struct { 
	Loglevel string
}

func readConfig(config_file_path string) ServerConfig {
	configFile, err := os.Open(config_file_path)
	if err != nil {
		fmt.Println(err)
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

func main_handler(w http.ResponseWriter, r *http.Request) {
}

func loglevel_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("----------------------------------------")
	fmt.Println("Received new request to change log level:")
	fmt.Println(r.Method, r.URL.Path)

	if !(r.Method == "GET" || r.Method == "PUT") {
		err := "Method = " + r.Method + " is not allowed to " + r.URL.Path
		fmt.Println(err)
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
			fmt.Println(res)
			http.Error(w, "400 bad request\n  Error: "+res, http.StatusBadRequest)
			return
		}

		if !contains_string(valid_log_levels, l.Loglevel) {
			slog.Error("Invalid log level")
			res := "Invalid log level " + l.Loglevel
			fmt.Println(res)
			http.Error(w, "400 bad request\n  Error: "+res, http.StatusBadRequest)
			return
		} else {
			app_log_level.Set(slog.LevelWarn)
			w.WriteHeader(http.StatusOK)
		}
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

	//fmt.Println(server_config)
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

	server_addr := server_config.Hostname + ":" + server_config.Port
	
	fmt.Printf("API server listening on: %s\n", server_addr)
	http.HandleFunc("/", main_handler)
	http.HandleFunc("/loglevel", loglevel_handler)
	http.ListenAndServe(server_addr, nil)
}