package middleware

import (
	"context"
	"file-transfer/pkg/common"
	"file-transfer/pkg/errno"
	"file-transfer/pkg/token"
	"net/http"
)

func AuthFilter(pf FiletransferHandlerFactory) FiletransferHandlerFactory {
	return func(next FiletransferContextHandlerFunc) FiletransferContextHandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			// log.C(ctx).Debugw("auth do start")
			idkey, userKey, err := token.ParseRequest(r)
			if err != nil {
				errno.WriteErrorResponse(ctx, w, errno.ErrAuthFail)
				return
			}
			ctx = context.WithValue(ctx, common.Trace_request_uid{}, idkey)
			ctx = context.WithValue(ctx, common.Trace_request_user{}, userKey)

			if pf != nil {
				next = pf(next)
			}
			next(ctx, w, r)
			// log.C(ctx).Debugw("auth do end")
		}
	}
}
