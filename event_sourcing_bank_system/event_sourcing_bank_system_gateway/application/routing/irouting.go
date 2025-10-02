package routing

import "event_sourcing_bank_system_gateway/application/model"

type RoutingUseCase interface {
	Forward(routingData *model.RoutingData) (interface{}, error)
}
