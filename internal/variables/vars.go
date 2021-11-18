package variables

import (
	"flag"
	"log"
)

var (
	TorAddr     = flag.String("tor", "localhost:9050", "tor address")
	TorCtlAddr  = flag.String("tor-ctl", "localhost:9051", "tor control address")
	Port        = flag.Uint("port", 2000, "gazer system server port")
	DSN         = flag.String("dsn", "./data.db", "DSN")
	TorPassword = flag.String("tor-pass", "password", "tor password")
	Concurrency = flag.Uint("concurrency", 8, "concurrency")
	inited      = false
)

func Init() {
	if !inited {
		flag.Parse()
		log.Println("TOR:             ", *TorAddr)
		log.Println("TOR_CTL:         ", *TorCtlAddr)
		log.Println("PORT:            ", *Port)
		log.Println("DSN:             ", *DSN)
		log.Println("CONCURRENCY:     ", *Concurrency)
		inited = true
	}
}
