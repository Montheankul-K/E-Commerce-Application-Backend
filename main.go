package main

import (
	"os"

	"github.com/Montheankul-K/E-Commerce-Application-Backend/config"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/servers"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/packages/databases"
)

func envPath() string {
	// argument : main args1 args2
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1] // .env.dev
	}
}

func main() {
	cfg := config.LoadConfig(envPath())

	db := databases.DbConnect(cfg.Db())
	defer db.Close() // defer ทำงานท้ายสุดก่อน func (main) จะหยุดทำงาน

	servers.NewServer(cfg, db).Start()
}
