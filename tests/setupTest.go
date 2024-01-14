package tests

import (
	"encoding/json"

	"github.com/Montheankul-K/E-Commerce-Application-Backend/config"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/servers"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/packages/databases"
)

func SetupTest() servers.IModuleFactory {
	cfg := config.LoadConfig("../.env.test")

	db := databases.DbConnect(cfg.Db())
	defer db.Close()

	s := servers.NewServer(cfg, db)
	return servers.InitModule(nil, s.GetServer(), nil)
}

func CompressToJSON(obj any) string {
	result, _ := json.Marshal(&obj)
	return string(result)
}
