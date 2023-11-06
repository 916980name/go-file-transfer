package controller

import (
	"context"
	"file-transfer/pkg/log"
	"net/http"
)

func Home(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	log.C(ctx).Infow("call Home")
	w.WriteHeader(200)
}
