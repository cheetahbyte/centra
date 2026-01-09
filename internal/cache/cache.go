package cache

import "sync/atomic"

var isReady atomic.Bool

func SetReady(v bool) {
	isReady.Store(v)
}

func IsReady() bool {
	return isReady.Load()
}
