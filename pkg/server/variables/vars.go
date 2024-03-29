package variables

import (
	"flag"
	"log"
)

var (
	TorAddr     = flag.String("tor", "localhost:9050", "tor address")
	TorCtlAddr  = flag.String("tor-ctl", "localhost:9051", "tor control address")
	SPEEDOS     = flag.String("speedos", "", "speedometer server's address")
	Port        = flag.Uint("port", 2000, "gazer system server port")
	DSN         = flag.String("dsn", "./data.db", "DSN")
	DB          = flag.String("db", "sqlite", "DB")
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
		log.Println("DB:              ", *DB)
		log.Println("CONCURRENCY:     ", *Concurrency)
		log.Println("SPEEDOS:         ", *SPEEDOS)
		inited = true
	}
}
