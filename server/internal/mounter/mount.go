package mounter

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"sync"
)

type MountHandler func(ctx context.Context, s *ghttp.Server)

var (
	funcOptions = make([]MountHandler, 0)
	fLock       sync.Mutex
)

func Mount(handler MountHandler) {
	fLock.Lock()
	defer fLock.Unlock()
	funcOptions = append(funcOptions, handler)
}

func DoMount(ctx context.Context, s *ghttp.Server) {
	for _, fn := range funcOptions {
		fn(ctx, s)
	}
}
