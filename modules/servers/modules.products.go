package servers

import (
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/products/productsHandlers"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/products/productsRepositories"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/products/productsUsecases"
)

type IProductsModule interface {
	Init()
	Repository() productsRepositories.IProductsRepository
	Usecase() productsUsecases.IProductsUsecase
	Handler() productsHandlers.IProductsHandler
}

type productsModule struct {
	*moduleFactory
	repository productsRepositories.IProductsRepository
	usecase    productsUsecases.IProductsUsecase
	handler    productsHandlers.IProductsHandler
}

func (m *moduleFactory) ProductsModule() IProductsModule {
	repository := productsRepositories.ProductsRepository(m.server.db, m.server.cfg, m.FilesModule().Usecase())
	usecase := productsUsecases.ProductsUsecase(repository)
	handler := productsHandlers.ProductsHandler(m.server.cfg, usecase, m.FilesModule().Usecase())

	return &productsModule{
		moduleFactory: m,
		repository:    repository,
		usecase:       usecase,
		handler:       handler,
	}
}

func (p *productsModule) Init() {
	router := p.router.Group("/products")
	router.Post("/", p.middleware.JwtAuth(), p.middleware.Authorize(2), p.handler.AddProduct)
	router.Patch("/:product_id", p.middleware.JwtAuth(), p.middleware.Authorize(2), p.handler.UpdateProduct)
	router.Get("/", p.middleware.ApiKeyAuth(), p.handler.FindProduct)
	router.Get("/:product_id", p.middleware.ApiKeyAuth(), p.handler.FindOneProduct)
	router.Delete("/:product_id", p.middleware.JwtAuth(), p.middleware.Authorize(2), p.handler.DeleteProduct)
}

func (p *productsModule) Repository() productsRepositories.IProductsRepository { return p.repository }
func (p *productsModule) Usecase() productsUsecases.IProductsUsecase           { return p.usecase }
func (p *productsModule) Handler() productsHandlers.IProductsHandler           { return p.handler }
