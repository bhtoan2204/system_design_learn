package service

import (
	"errors"
	"event_sourcing_bank_system_gateway/application/model"
	"event_sourcing_bank_system_gateway/package/settings"
)

type serviceClient struct {
	remoteServiceClientRegistry map[string]RemoteServiceClient
}

func NewServiceClient(config *settings.Config) *serviceClient {
	return &serviceClient{
		initRemoteServiceClientRegistry(config),
	}
}

func (service *serviceClient) Invoke(routingData *model.RoutingData) (interface{}, error) {
	remoteServiceClient, found := service.remoteServiceClientRegistry[routingData.ServiceName]
	if !found {
		return nil, errors.New("remote service client not found: " + routingData.ServiceName)
	}
	return remoteServiceClient.Invoke(routingData.ServiceMethod, routingData.Payload, routingData.Metadata)
}
