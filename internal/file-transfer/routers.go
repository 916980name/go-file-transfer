package filetransfer

import (
	"context"
	"net/http"

	"file-transfer/internal/file-transfer/controller"
	"file-transfer/internal/file-transfer/repo"
	"file-transfer/internal/file-transfer/service"
	"file-transfer/pkg/db/dbmongo"
	"file-transfer/pkg/middleware"

	"github.com/gorilla/mux"
)

func handleMuxChain(f middleware.FiletransferContextHandlerFunc) middleware.FiletransferHandlerFactory {
	return func(next middleware.FiletransferContextHandlerFunc) middleware.FiletransferContextHandlerFunc {
		return f
	}
}

func handleMuxChainFunc(pf middleware.FiletransferHandlerFactory) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f := pf(nil)
		f(nil, w, r)
	}
}

func wrapper(f middleware.FiletransferContextHandlerFunc) http.HandlerFunc {
	return handleMuxChainFunc(middleware.RequestFilter(handleMuxChain(f)))
}

func initAllRouters(r *mux.Router) error {
	messageRepo := repo.NewMessageRepo(dbmongo.GetClient(context.Background()))
	messageService := service.NewMessageService(messageRepo)
	messageController := controller.NewMessageController(messageService)

	r.NewRoute().Methods("GET").Path("/home").HandlerFunc(wrapper(controller.Home))
	r.NewRoute().Methods("GET").Path("/msg").HandlerFunc(wrapper(messageController.Onemessage))
	return nil
}
