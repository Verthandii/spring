package main

import (
	"flag"
	"log"
)

var (
	root = flag.String("root", "api", "root path")
)

func main() {
	flag.Parse()
	err := generateFile(*root)
	if err != nil {
		log.Println(err)
	}
}
