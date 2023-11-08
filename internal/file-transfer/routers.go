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

func authWrapper(f middleware.FiletransferContextHandlerFunc) http.HandlerFunc {
	return handleMuxChainFunc(middleware.RequestFilter(middleware.AuthFilter(handleMuxChain(f))))
}

func initAllRouters(r *mux.Router) error {
	mongoClient := dbmongo.GetClient(context.Background())
	messageRepo := repo.NewMessageRepo(mongoClient)
	messageService := service.NewMessageService(messageRepo)
	messageController := controller.NewMessageController(messageService)
	userRepo := repo.NewUserRepo(mongoClient)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	r.NewRoute().Methods("GET").Path("/home").HandlerFunc(wrapper(controller.Home))
	r.NewRoute().Methods("POST").Path("/trysignin").HandlerFunc(wrapper(userController.Login))

	r.NewRoute().Methods("POST").Path("/msg").HandlerFunc(authWrapper(messageController.ReadMessage))
	r.NewRoute().Methods("PUT").Path("/msg").HandlerFunc(authWrapper(messageController.SendMessage))
	r.NewRoute().Methods("DELETE").Path("/msg/{mId}").HandlerFunc(authWrapper(messageController.DeleteMessage))
	return nil
}
