package servers

import (
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/files/filesHandlers"
	"github.com/Montheankul-K/E-Commerce-Application-Backend/modules/files/filesUsecases"
)

type IFilesModule interface {
	Init()
	Usecase() filesUsecases.IFilesUsecase
	Handler() filesHandlers.IFilesHandler
}

type fileModule struct {
	*moduleFactory
	usecase filesUsecases.IFilesUsecase
	handler filesHandlers.IFilesHandler
}

func (m *moduleFactory) FilesModule() IFilesModule {
	usecase := filesUsecases.FilesUsecase(m.server.cfg)
	handler := filesHandlers.FilesHandler(m.server.cfg, usecase)

	return &fileModule{
		moduleFactory: m,
		usecase:       usecase,
		handler:       handler,
	}
}

func (f *fileModule) Init() {
	router := f.router.Group("files")
	router.Post("/upload", f.middleware.JwtAuth(), f.middleware.Authorize(2), f.handler.UploadFiles)
	router.Patch("/delete", f.middleware.JwtAuth(), f.middleware.Authorize(2), f.handler.DeleteFile)
}

func (f *fileModule) Usecase() filesUsecases.IFilesUsecase { return f.usecase }
func (f *fileModule) Handler() filesHandlers.IFilesHandler { return f.handler }
