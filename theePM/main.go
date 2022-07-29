package main

import (
	"log"
	"theePM/paramMiner"
)

func main() {
	opt := paramMiner.ParseFlags()
	err := opt.Start()
	if err != nil{
		log.Fatal(err)
	}

}
