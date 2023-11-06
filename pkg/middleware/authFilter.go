package middleware

import (
	"context"
	"file-transfer/pkg/log"
	"net/http"
)

func AuthFilter(pf FiletransferHandlerFactory) FiletransferHandlerFactory {
	return func(next FiletransferContextHandlerFunc) FiletransferContextHandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			log.C(ctx).Infow("auth do start")
			if pf != nil {
				next = pf(next)
			}
			next(ctx, w, r)
			log.C(ctx).Infow("auth do end")
		}
	}
}
