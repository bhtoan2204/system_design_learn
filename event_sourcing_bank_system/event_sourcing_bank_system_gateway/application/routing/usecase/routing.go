package usecase

import (
	"event_sourcing_bank_system_gateway/application/model"
	"event_sourcing_bank_system_gateway/application/routing"
)

type routingUseCase struct {
	serviceClient routing.ServiceClient
}

func NewRoutingUseCase(serviceClient routing.ServiceClient) *routingUseCase {
	return &routingUseCase{
		serviceClient: serviceClient,
	}
}

func (rUC *routingUseCase) Forward(routingData *model.RoutingData) (interface{}, error) {
	return rUC.serviceClient.Invoke(routingData)
}
