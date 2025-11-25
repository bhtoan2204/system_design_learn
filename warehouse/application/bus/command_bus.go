package bus

import (
	"context"
	"fmt"
	"warehouse/pkg/command"
)

type bus struct {
	handlers    map[string]command.HandlerFunc
	middlewares []command.MiddlewareFunc
}

func NewCommandBus() command.CommandBus {
	return &bus{
		handlers:    make(map[string]command.HandlerFunc),
		middlewares: make([]command.MiddlewareFunc, 0),
	}
}

func (b *bus) Register(ctx context.Context, command command.Command, handler command.HandlerFunc) error {
	if err := b.validate(command); err != nil {
		return err
	}
	b.handlers[command.CommandName()] = handler
	return nil
}

func (b *bus) UseMiddleware(ctx context.Context, middleware command.MiddlewareFunc) error {
	b.middlewares = append(b.middlewares, middleware)
	return nil
}

func (b *bus) Dispatch(ctx context.Context, command command.Command) error {
	handler, ok := b.handlers[command.CommandName()]
	if !ok {
		return fmt.Errorf("command not found: %s", command.CommandName())
	}
	for _, middleware := range b.middlewares {
		handler = middleware(handler)
	}
	return handler(ctx, command)
}
