package utils

import (
	"context"
	"net/http"
)

type CtxValues struct {
	Method    string
	Url       string
	Agent     string
	UserID    string
	RemoteIP  string
	RequestID string
}
type ctxKeyType string

var ctxValueKey = ctxKeyType("_vFER__")

func CtxGenerate(req *http.Request, userID string, requestID string) context.Context {
	return context.WithValue(req.Context(), ctxValueKey, &CtxValues{
		Method:    req.Method,
		Url:       req.URL.Path,
		Agent:     req.UserAgent(),
		UserID:    userID,
		RemoteIP:  req.Header.Get("X-Forwarded-For"),
		RequestID: requestID,
	})
}
