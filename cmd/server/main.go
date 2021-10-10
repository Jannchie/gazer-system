package main

import "github.com/jannchie/gs-server/internal/gs"

func main() {
	gss := gs.NewDefaultServer()
	gss.Run(2000)
}
