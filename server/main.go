package main

import (
)

var config Config

func main() {
	config.Load()
	ApiServer()
}
