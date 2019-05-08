package main

import (
	"log"
	"net/http"
	"encoding/json"
	"text/template"
)

var errJSON = func(msg string) []byte {
	return []byte("{\"error\": \"" + msg + "\"}")
}

func Render(obj interface{}, err error) []byte {
	if err != nil {
		return errJSON(err.Error())
	}
	payload, err := json.Marshal(obj)
	if err != nil {
		return errJSON(err.Error())
	}
	return payload
}

func index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./index.html")
	t.Execute(w, map[string] string {"Prefix": config.Prefix})
}

func listFiles(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()["path"]
	path := "/"
	if len(params) != 0 {
		path = params[0]
	}
	w.Write(Render(GetAllFiles(path)))
}

func ApiServer()  {
	http.HandleFunc("/", M(index, Logging))
	http.HandleFunc("/list", M(listFiles, CORS, Logging))
	log.Printf("Listening on localhost:8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
