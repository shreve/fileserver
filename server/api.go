package main

import (
	"io/ioutil"
	"fmt"
	"log"
	"path"
	"errors"
	"net/url"
	"net/http"
	"encoding/json"
	"text/template"
)

var errJSON = func(msg string) []byte {
	return []byte("{\"error\": \"" + msg + "\"}")
}

var PathRequired = errors.New("You must supply a path parameter with that request")

type Params struct {
	Path string  `json:"path"`
}

func Render(obj interface{}, err error) []byte {
	if err != nil {
		return errJSON(err.Error())
	}
	payload, err := json.Marshal(obj)
	if err != nil {
		log.Printf("Returning error: " + err.Error())
		return errJSON(err.Error())
	}
	return payload
}

func index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./index.html")
	t.Execute(w, config)
}

func listFiles(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()["path"]
	path := "/"
	if len(params) != 0 {
		path = params[0]
	}
	w.Write(Render(GetAllFiles(path)))
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	reqpath := pathParam(r)
	if reqpath != "" {
		file, err := DownloadFile(reqpath)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(Render(nil, err))
			return
		}
		attachment := fmt.Sprintf("attachment; filename=\"%s\"", path.Base(file))
		w.Header().Set("Content-Disposition", attachment)
		w.Header().Set("Content-Type", mimeType(file))
		w.Header().Set("Location", "/download?path=" + reqpath)
		log.Printf("Sending file: %s", file)
		http.ServeFile(w, r, file)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(Render(nil, PathRequired))
	}
}

func statusFile(w http.ResponseWriter, r *http.Request) {
	reqpath := pathParam(r)
	if reqpath != "" {
		err := StatusFile(reqpath)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(Render(nil, err))
			return
		}
		w.Write(Render(nil, nil))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(Render(nil, PathRequired))
		return
	}
}

func mimeType(file string) string {
	switch ext := path.Ext(file); ext {
	case ".pdf": return "application/pdf"
	case ".zip": return "application/zip"
	case ".mp3": return "audio/mpeg"
	default:
		log.Printf(ext)
		return "application/octet-stream"
	}
}

func pathParam(r *http.Request) string {
	path := ""
	if r.Method == "GET" {
		params := r.URL.Query()["path"]
		if len(params) != 0 {
			path, _ = url.PathUnescape(params[0])
		}
	}
	if r.Method == "POST" {
		buf, _ := ioutil.ReadAll(r.Body)
		params := Params{}
		json.Unmarshal(buf, &params)
		path = params.Path
	}
	return path
}

func ApiServer()  {
	http.HandleFunc("/", M(index, Logging))
	http.HandleFunc("/api", M(index, Logging))
	http.HandleFunc("/list", M(listFiles, CORS, Logging))
	http.HandleFunc("/api/list", M(listFiles, CORS, Logging))
	http.HandleFunc("/download", M(downloadFile, CORS, Logging))
	http.HandleFunc("/api/download", M(downloadFile, CORS, Logging))
	http.HandleFunc("/status", M(statusFile, CORS, Logging))
	http.HandleFunc("/api/status", M(statusFile, CORS, Logging))
	log.Printf("Listening on localhost:8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
