package middleware

import (
	"context"
	"file-transfer/pkg/common"
	"file-transfer/pkg/log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type FiletransferHandlerFunc func(w http.ResponseWriter, r *http.Request)
type FiletransferContextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request)
type FiletransferHandlerFactory func(FiletransferContextHandlerFunc) FiletransferContextHandlerFunc

func RequestFilter(pf FiletransferHandlerFactory) FiletransferHandlerFactory {
	return func(next FiletransferContextHandlerFunc) FiletransferContextHandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			if ctx == nil {
				ctx = context.Background()
			}
			ip := getClientIP(r)
			ctx = context.WithValue(ctx, common.RESOURCE_IP, ip)

			rId := getRequestId(r)
			ctx = context.WithValue(ctx, common.REQUEST_ID, rId)

			ctx = context.WithValue(ctx, common.REQUEST_URI, getRequestUri(r))
			ctx = context.WithValue(ctx, common.REQUEST_METHOD, getRequestMethod(r))
			ctx = context.WithValue(ctx, common.REQUEST_USER, getRequestUser(r))
			log.C(ctx).Debugw("new Request --> ")
			if pf != nil {
				next = pf(next)
			}
			next(ctx, w, r)
		}
	}
}

func getRequestUri(r *http.Request) string {
	return r.URL.Path
}

func getRequestMethod(r *http.Request) string {
	return r.Method
}

func getRequestUser(r *http.Request) string {
	return r.Header.Get(common.REQUEST_USER)
}

func getRequestId(r *http.Request) string {
	requestID := r.Header.Get(common.REQUEST_ID)
	if requestID == "" {
		requestID = uuid.New().String()
	}
	return requestID
}

func getClientIP(r *http.Request) string {
	// Check for X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// The client IP may be a comma-separated list, get the first IP
		ips := strings.Split(xff, ",")
		ip := strings.TrimSpace(ips[0])
		return ip
	}

	// Check for X-Real-IP header
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}

	// If headers not found, use RemoteAddr as fallback
	ip := r.RemoteAddr

	// If the IP address contains a port number, remove it
	if index := strings.IndexByte(ip, ':'); index >= 0 {
		ip = ip[:index]
	}

	return ip
}
