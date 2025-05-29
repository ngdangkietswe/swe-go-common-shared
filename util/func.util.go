package util

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
)

// HandleGrpc safely executes a gRPC handler function with panic recovery.
func HandleGrpc[Req any, Resp any](ctx context.Context, req Req, handler func(context.Context, Req) (Resp, error)) (resp Resp, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[PANIC RECOVERED] %v\nStack: %s", r, string(debug.Stack()))
			err = fmt.Errorf("internal server error")
		}
	}()

	return handler(ctx, req)
}
