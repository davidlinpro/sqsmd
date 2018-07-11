package main

import "log"

func init() {
	Config.Init()
	Cache.Init()
}

var ()

func main() {
	log.Println("starting server...")
	go processQueueStats()
	handleRequests()
}
