package mlog

import (
	"context"

	"google.golang.org/grpc/metadata"
)

const (
	//ContextKeyRequestID reqid
	ContextKeyRequestID = "reqid"
)

//GetRequestID 获取RequestID
func GetRequestID(ctx context.Context) (has bool, reqID string) {
	md, has := metadata.FromIncomingContext(ctx)
	if !has {
		return
	}
	val, has := md[ContextKeyRequestID]
	if has {
		reqID = val[0]
	}
	return
}
