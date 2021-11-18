package main

import (
	"log"

	"github.com/jannchie/gazer-system/internal/gs"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	gss := gs.NewDefaultServer()
	gss.Run()
}
