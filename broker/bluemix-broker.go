package broker

import (
	"errors"
	"github.com/bcshuai/brokerapi"
	"github.com/bcshuai/cf-redis-broker/brokerconfig"
)

type BluemixServiceMetadata struct {
	brokerapi.ServiceMetadata
	ServiceKeysSupported 	bool 	`json:"serviceKeysSupported,omitempty"`
	Type 					string  `json:"type, omitempty"`
}

type BluemixServicePlan struct {
	brokerapi.ServicePlan
	MaxMemoryInMB 			int 		`json:"max_memory_mb"`
	MaxClientConnections	int 		`json:"max_client_connection"`
}

type BluemixRedisServiceBroker struct {
	RedisServiceBroker
}

func (broker *BluemixRedisServiceBroker) Services() []brokerapi.IMetadataProvider {
	config := broker.Config.RedisConfiguration
	services := []brokerapi.IMetadataProvider{}
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

func (broker *BluemixRedisServiceBroker) Provision(instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error){
	
	spec := brokerapi.ProvisionedServiceSpec{
		false,
		"",
	}

	if broker.instanceExists(instanceID) {
		return spec, brokerapi.ErrInstanceAlreadyExists
	}

	if details.PlanID == "" {
		return spec, errors.New("plan_id required")
	}

	var matchedServiceConfig brokerconfig.BluemixServiceConfig

	for _, serviceConfig := range broker.Config.RedisConfiguration.Services {
		if(serviceConfig.ServiceID == details.ServiceID) {
			matchedServiceConfig = serviceConfig
			break
		}
		return spec, errors.New("no such service")
	}

	var matchedServicePlanConfig brokerconfig.BluemixServicePlanConfig
	for _, planConfig := range matchedServiceConfig.Plans {
		if planConfig.ID == details.PlanID {
			matchedServicePlanConfig = planConfig
			break
		}
		return spec, errors.New("no such plan: " + details.PlanID)
	}

	instanceCreatorKey := ""
	if(matchedServiceConfig.ServiceName == "redis-shared-vm"){
		instanceCreatorKey = "shared"
	} else if (matchedServiceConfig.ServiceName == "dedicated") {
		instanceCreatorKey = "dedicated"
	}

	if instanceCreatorKey == "" {
		return spec, errors.New("Service ID not recognized")
	}

	instanceCreator, ok := broker.InstanceCreators[instanceCreatorKey]
	if !ok {
		return spec, errors.New("instance creator not found for plan")
	}

	return spec, instanceCreator.CreateWithRestriction(instanceID,matchedServicePlanConfig.MaxMemoryInMB,matchedServicePlanConfig.MaxClientConnections)
}

func getServiceMetaFromConfig(serviceMetadataConfig brokerconfig.BluemixServiceMetadataConfig) BluemixServiceMetadata {
	return BluemixServiceMetadata {
		ServiceMetadata: brokerapi.ServiceMetadata {
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
func getServicePlansFromConfigs(servicePlanConfigs []brokerconfig.BluemixServicePlanConfig) []brokerapi.IMetadataProvider {
	plans := []brokerapi.IMetadataProvider{}
	for _, planConfig := range servicePlanConfigs {
		plan := BluemixServicePlan {
			ServicePlan: brokerapi.ServicePlan {
				ID:          planConfig.ID,
				Name:        planConfig.Name,
				Description: planConfig.Description,
				Metadata: getPlanMetaFromConfig(planConfig.Metadata), 
			},
			MaxMemoryInMB: planConfig.MaxMemoryInMB,
			MaxClientConnections: planConfig.MaxClientConnections, 
		}
		plans = append(plans, plan)
	}
	return plans
}
func getPlanMetaFromConfig(servicePlanMetadataConfig brokerconfig.BluemixServicePlanMetadataConfig) brokerapi.ServicePlanMetadata {
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


