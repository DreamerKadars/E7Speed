package main

import (
	"E7Speed/httpserver"
	"E7Speed/model"
)

func init() {
	// go service.DealOldFile()
    model.InitModel()
}

func main() {
	httpserver.StartHttpServer()
}
