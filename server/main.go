package main

import (
)

var config Config

func main() {
	config.Load()
	config.Inspect()
	ApiServer()
	// port := flag.String("p", "3000", "port")
	// dir := flag.String("d", ".", "dir")
	// flag.Parse()

	// http.Handle("/", http.FileServer(http.Dir(*dir)))
	// log.Printf("Serving %s on http port %s\n", *dir, *port)
	// log.Fatal(http.ListenAndServe(":"+*port, nil))
}
