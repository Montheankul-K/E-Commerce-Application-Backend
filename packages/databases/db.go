// migrations : initialize database
package databases

import (
	"log"

	"github.com/Montheankul-K/E-Commerce-Application-Backend/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func DbConnect(cfg config.IDbConfig) *sqlx.DB {
	// หรือ config.IConfig ดึงมาทั้งหมดก็ได้
	db, err := sqlx.Connect("pgx", cfg.Url()) // sqlx จะไปเรียกใช้ driver pgx
	if err != nil {
		log.Fatalf("Connect to db failed: %v", err)
	}
	db.DB.SetMaxOpenConns(cfg.MaxOpenConnections())
	return db
}
