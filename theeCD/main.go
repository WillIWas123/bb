package main

import (
	"log"
	"theeCD/contentDiscovery"
)

func main() {
	opt := contentDiscovery.ParseFlags()
	err := opt.Start()
	if err != nil{
		log.Fatal(err)
	}

}
