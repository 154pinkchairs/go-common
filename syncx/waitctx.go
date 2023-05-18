package syncx

import (
	"context"
	"sync"
	"time"
)

type WaitCtx struct {
	ctx context.Context
	wg  sync.WaitGroup
}

func (wc *WaitCtx) Add(ctx context.Context, delta int) {
	wc.wg.Add(delta)
	go func() {
		select {
		case <-ctx.Done():
			wc.wg.Add(-delta)
		}
	}()
}

func (wc *WaitCtx) AddTimeout(delta int, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(wc.ctx, timeout)
	defer cancel()
	wc.Add(ctx, delta)
}

func (wc *WaitCtx) Done() <-chan struct{} {
	return wc.ctx.Done()
}

func (wc *WaitCtx) Wait() {
	wc.wg.Wait()
}

func (wc *WaitCtx) WaitTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(wc.ctx, timeout)
	defer cancel()
	return wc.WaitContext(ctx)
}

func (wc *WaitCtx) WaitContext(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		wc.wg.Wait()
		close(done)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}
