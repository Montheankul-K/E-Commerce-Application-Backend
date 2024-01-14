// Modules : internal services eg. api, use case, repository, handler, ..
package servers

import (
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/appinfo/addinfoHandlers"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/appinfo/appinfoRepositories"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/appinfo/appinfoUsecases"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/files/filesUsecases"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/middlewares/middlewaresHandlers"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/middlewares/middlewaresRepositories"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/middlewares/middlewaresUsecases"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/monitor/monitorHandlers"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/orders/ordersHandlers"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/orders/ordersRepositories"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/orders/ordersUsecases"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/products/productsRepositories"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/users/usersHandlers"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/users/usersRepositories"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppinfoModule()
	FilesModule() IFilesModule
	ProductsModule() IProductsModule
	OrdersModule()
}

type moduleFactory struct {
	router     fiber.Router
	server     *server
	middleware middlewaresHandlers.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, mid middlewaresHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		router:     r,
		server:     s,
		middleware: mid,
	}
}

func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewaresHandler {
	repository := middlewaresRepositories.MiddlewaresRepository(s.db)
	usecase := middlewaresUsecases.MiddlewaresUsecase(repository)
	return middlewaresHandlers.MiddlewaresHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.server.cfg)

	m.router.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UsersRepository(m.server.db)
	usecase := usersUsecases.UsersUsecase(m.server.cfg, repository)
	handler := usersHandlers.UsersHandler(m.server.cfg, usecase)

	router := m.router.Group("/users")
	router.Post("/signup", m.middleware.ApiKeyAuth(), handler.SignUpCustomer)
	router.Post("/signin", m.middleware.ApiKeyAuth(), handler.SignIn)
	router.Post("/refresh", m.middleware.ApiKeyAuth(), handler.RefreshPassport)
	router.Post("/signout", m.middleware.ApiKeyAuth(), handler.SignOut)
	router.Post("/signup-admin", m.middleware.JwtAuth(), m.middleware.Authorize(2), handler.SignOut)
	router.Get("/admin/secret", m.middleware.JwtAuth(), m.middleware.Authorize(2), handler.GenerateAdminToken)
	// initial admin (sql migration) > generate admin key > ส่ง admin token ผ่าน middlewares ทุกครั้งที่ signup admin
	router.Get("/:user_id", m.middleware.JwtAuth(), m.middleware.ParamsCheck(), handler.GetUserProfile)
}

func (m *moduleFactory) AppinfoModule() {
	repository := appinfoRepositories.AppinfoRepository(m.server.db)
	usecase := appinfoUsecases.AppinfoUsecase(repository)
	handler := appinfoHandlers.AppinfoHandler(m.server.cfg, usecase)

	router := m.router.Group("/appinfo")
	router.Post("/categories", m.middleware.JwtAuth(), m.middleware.Authorize(2), handler.AddCategory)
	router.Get("/categories", m.middleware.ApiKeyAuth(), handler.FindCategory)
	router.Get("/apikey", m.middleware.JwtAuth(), m.middleware.Authorize(2), handler.GenerateApiKey)
	router.Delete("/:category_id/categories", m.middleware.JwtAuth(), m.middleware.Authorize(2), handler.RemoveCategory)
}

func (m *moduleFactory) OrdersModule() {
	filesUsecase := filesUsecases.FilesUsecase(m.server.cfg)
	productsRepository := productsRepositories.ProductsRepository(m.server.db, m.server.cfg, filesUsecase)

	repository := ordersRepositories.OrdersRepository(m.server.db)
	usecase := ordersUsecases.OrdersUsecase(repository, productsRepository)
	handler := ordersHandlers.OrdersHandler(m.server.cfg, usecase)

	router := m.router.Group("/orders")
	router.Post("/", m.middleware.JwtAuth(), handler.InsertOrder)
	router.Get("/", m.middleware.JwtAuth(), m.middleware.Authorize(2), handler.FindOrder)
	router.Get("/:user_id/:order_id", m.middleware.JwtAuth(), m.middleware.ParamsCheck(), handler.FindOneOrder)
	router.Patch("/:user_id/:order_id", m.middleware.JwtAuth(), m.middleware.ParamsCheck(), handler.UpdateOrder)
}
