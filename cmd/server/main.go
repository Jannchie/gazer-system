package main

import (
	"github.com/jannchie/gazer-system/internal/gs"
)

func main() {
	gss := gs.NewDefaultServer()
	gss.Run()
}
