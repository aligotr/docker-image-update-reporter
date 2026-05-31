package registry

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.podman.io/image/v5/types"
)

type Client struct {
	opts   Options
	sysCtx *types.SystemContext
}

type Options struct {
	Timeout time.Duration
}

// Новый инстанс docker-registry
func New() (*Client, error) {
	opts := Options{
		Timeout: 20 * time.Second,
	}

	return &Client{
		opts:   opts,
		sysCtx: &types.SystemContext{},
	}, nil
}

func (c *Client) timeoutContext() (context.Context, context.CancelFunc) {
	ctx := context.Background()
	var cancelFunc context.CancelFunc = func() {}
	if c.opts.Timeout > 0 {
		cancelCtx, cancelCause := context.WithCancelCause(ctx)

		var cancelTimeout context.CancelFunc
		ctx, cancelTimeout = context.WithTimeoutCause(cancelCtx, c.opts.Timeout, errors.New(fmt.Sprint(context.DeadlineExceeded)))

		cancelFunc = func() {
			cancelTimeout()
			cancelCause(errors.New(fmt.Sprint(context.Canceled)))
		}
	}
	return ctx, cancelFunc
}
