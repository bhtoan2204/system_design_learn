package handler

import "event_sourcing_bank_system_gateway/package/wrapper"

type RequestHandler interface {
	Handle(ctx *wrapper.Context) (interface{}, error)
}
