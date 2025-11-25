package command

import (
	"context"
)

type HandlerFunc func(ctx context.Context, command Command) error
type MiddlewareFunc func(h HandlerFunc) HandlerFunc

type CommandBus interface {
	Register(ctx context.Context, command Command, handler HandlerFunc) error
	UseMiddleware(ctx context.Context, middleware MiddlewareFunc) error
	Dispatch(ctx context.Context, command Command) error
}
