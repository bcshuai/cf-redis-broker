package broker

import (
	"github.com/bcshuai/brokerapi"
)

type BluemixServiceMetadata struct {
	brokerapi.ServiceMetadata
	ServiceKeysSupported 	bool 	`json:"serviceKeysSupported,omitempty"`
	Type 					string  `json:"type, omitempty"`
}

type BluemixServicePlan struct {
	brokerapi.ServicePlan
	MaxMemoryInMB 			int 		`json:"max_memory_mb"`
	MaxClientConnections	int 		`json:"max_user_connection"`
}

type BluemixRedisServiceBroker struct {
	RedisServiceBroker
}

func (broker *BluemixRedisServiceBroker) Services() []brokerapi.IMetadataProvider {
	config := broker.Config
	services := []IMetadataProvider{}
	for _, serviceConfig := range config.Services {
		service := brokerapi.Service {
			ID:          serviceConfig.ServiceID,
			Name:        serviceConfig.ServiceName,
			Description: serviceConfig.Description,
			Bindable:    true,
			Plans:       getServicePlansFromConfigs(serviceConfig.Plans),
			Metadata: getServiceMetaFromConfig(serviceConfig.Metadata),
			Tags: serviceConfig.Tags,
		}
		services = append(services, service)
	}
	return services
}
func getServiceMetaFromConfig(serviceMetadataConfig BluemixServiceMetadataConfig) BluemixServiceMetadata {
	return BluemixServiceMetadata {
		brokerapi.ServiceMetadata {
			DisplayName:      serviceMetadataConfig.DisplayName,
			LongDescription:  serviceMetadataConfig.LongDescription,
			DocumentationUrl: serviceMetadataConfig.DocumentationUrl,
			SupportUrl:       serviceMetadataConfig.SupportUrl,
			ImageUrl: 		  serviceMetadataConfig.ImageUrl,
			ProviderDisplayName: serviceMetadataConfig.ProviderDisplayName,
		},
		ServiceKeysSupported: serviceMetadataConfig.ServiceKeysSupported,
		Type: serviceMetadataConfig.Type,
	}
}
func getServicePlansFromConfigs(servicePlanConfigs []BluemixServicePlanConfig) []BluemixServicePlan {
	plans := []BluemixServicePlan{}
	for _, planConfig := range servicePlanConfigs {
		plan := BluemixServicePlan {
			broker.ServicePlan {
				ID:          planConfig.ID,
				Name:        planConfig.Name,
				Description: planConfig.Description,
				Metadata: getPlanMetaFromConfig(planConfig.Metadata), 
			},
			MaxMemoryInMB: planConfig.MaxMemoryInMB,
			MaxClientConnections: planConfig.MaxClientConnections, 
		}
	}
}
func getPlanMetaFromConfig(servicePlanMetadataConfig BluemixServicePlanMetadataConfig) brokerapi.ServicePlanMetadata {
	costs := []brokerapi.ServiceCost{}
	costConfigs := servicePlanMetadataConfig.Costs
	for _, costConfig := range costConfigs {
		cost := brokerapi.ServiceCost {
			Amount: costConfig.Amount,
			Unit: costConfig.Unit,
		}
		costs = append(costs, cost)
	}
	return brokerapi.ServicePlanMetadata {
		DisplayName: servicePlanMetadataConfig.DisplayName,
		Bullets: servicePlanMetadataConfig.Bullets,
		Costs: costs,
	}
}


