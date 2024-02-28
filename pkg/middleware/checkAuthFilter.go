package middleware

import (
	"context"
	"file-transfer/pkg/common"
	"file-transfer/pkg/errno"
	"file-transfer/pkg/log"
	"net/http"
)

func CheckAuthFilter(pf FiletransferHandlerFactory) FiletransferHandlerFactory {
	return func(next FiletransferContextHandlerFunc) FiletransferContextHandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			log.C(ctx).Debugw("check auth do start")
			if user := ctx.Value(common.CTX_USER_NAME); user == nil {
				errno.WriteErrorResponse(ctx, w, errno.ErrAuthFail)
				return
			}
			if uid := ctx.Value(common.CTX_USER_KEY); uid == nil {
				errno.WriteErrorResponse(ctx, w, errno.ErrAuthFail)
				return
			}
			if pf != nil {
				next = pf(next)
			}
			next(ctx, w, r)
			log.C(ctx).Debugw("check auth do end")
		}
	}
}
