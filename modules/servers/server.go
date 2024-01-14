package servers

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/Montheankul-K/E-Commerce-Application-Backend/config"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type IServer interface {
	Start()
	GetServer() *server
}

type server struct {
	app *fiber.App
	cfg config.IConfig
	db  *sqlx.DB
}

func NewServer(cfg config.IConfig, db *sqlx.DB) IServer {
	return &server{
		cfg: cfg,
		db:  db,
		app: fiber.New(fiber.Config{
			AppName:      cfg.App().Name(),
			BodyLimit:    cfg.App().BodyLimit(),
			ReadTimeout:  cfg.App().ReadTimeout(),
			WriteTimeout: cfg.App().WriteTimeout(),
			JSONEncoder:  json.Marshal,
			JSONDecoder:  json.Unmarshal,
		}),
	}
}

func (s *server) GetServer() *server {
	return s
}

func (s *server) Start() {
	// middlewares
	middlewares := InitMiddlewares(s)
	s.app.Use(middlewares.Logger())
	s.app.Use(middlewares.Cors()) // ประกาศให้ middlewares เป็น global สำหรับ end point ใดๆ (เข้า middlewares ก่อนทุก end point)

	// modules
	v1 := s.app.Group("v1") // เพิ่่ม prefix v1 : https://localhost:3000/v1
	modules := InitModule(v1, s, middlewares)
	modules.MonitorModule()
	modules.UsersModule()
	modules.AppinfoModule()
	modules.FilesModule().Init()
	modules.ProductsModule().Init()
	modules.OrdersModule()

	s.app.Use(middlewares.RouterCheck())

	// gaceful shutdown : คืน resource ทั้งหมด (ค่อยๆ shutdown) ถ้า server ถูก interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		log.Println("server is shutting down...")
		_ = s.app.Shutdown()
	}()

	// listen to host:port
	log.Printf("server is starting on %v", s.cfg.App().Url())
	s.app.Listen(s.cfg.App().Url())
}
