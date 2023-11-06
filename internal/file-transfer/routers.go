package filetransfer

import (
	"net/http"

	"file-transfer/internal/file-transfer/controller"
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

func initAllRouters(r *mux.Router) error {
	r.NewRoute().Methods("GET").Path("/home").HandlerFunc(handleMuxChainFunc(middleware.RequestFilter(handleMuxChain(controller.Home))))
	return nil
}
