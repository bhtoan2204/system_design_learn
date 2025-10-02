package routing

import "event_sourcing_bank_system_gateway/application/model"

type ServiceClient interface {
	Invoke(routingData *model.RoutingData) (interface{}, error)
}
