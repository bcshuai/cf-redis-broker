package agent-sharednode

import (
	"errors"

	"github.com/pborman/uuid"

	"github.com/bcshuai/brokerapi"
	"github.com/bcshuai/cf-redis-broker/redis"
	sharenode "github.com/bcshuai/cf-redis-broker/agentapi-sharednode"
	"github.com/bcshuai/cf-redis-broker/brokerconfig"
	"github.com/bcshuai/cf-redis-broker/broker"
)

type SharedNodeClient struct {
	InstanceCreator broker.InstanceCreator
	InstanceBinder  broker.InstanceBinder
	Config           brokerconfig.Config
}

func (client *SharedNodeClient) Resources() (Resource, error) {
	limitation := client.Config.RedisConfiguration.ServiceInstanceLimit
	used := client.InstanceCreator.InstanceCount()
	instanceStatus := sharenode.ResourceStatus{
		All: limitation,
		Used: used,
		Free: limitation - used,
	}

	return sharenode.Resource{
		InstanceStatus: instanceStatus,
	}, nil
}

func (client *SharedNodeClient) AllInstances() ([]*redis.Instance, error) {
	return client.InstanceCreator.AllInstances()
}

func (client *SharedNodeClient) InstanceInfo(instanceId string) (*redis.Instance, error) {
	return client.InstanceCreator.FindByID(instanceId)
}

func (client *SharedNodeClient) InstanceExists(instanceId string) (bool, error) {
	return client.InstanceCreator.InstanceExists(instanceId)
}

func (client *SharedNodeClient) InstanceCredential(instanceId string) (broker.InstanceCredentials, error) {
	ok, err := client.InstanceCreator.InstanceExists(instanceId)
	if(err != nil){
		return broker.InstanceCredentials{}, err
	}
	if !ok {
		return sharenode.RedisCredential{}, brokerapi.ErrInstanceDoesNotExist
	}
	instance, err := client.InstanceCreator.FindByID(instanceId)
	if(err != nil){
		return broker.InstanceCredentials{}, err
	}

	return broker.InstanceCredentials{
		Host:     instance.Host,
		Port:     instance.Port,
		Password: instance.Password,
	}, nil
}

func (client *SharedNodeClient) ProvisionInstance(instance redis.Instance) error {
	instanceCount, err := client.InstanceCreator.InstanceCount()
	if err != nil {
		return err
	}

	if instanceCount >= client.InstanceCreator.RedisConfiguration.ServiceInstanceLimit {
		return brokerapi.ErrInstanceLimitMet
	}

	instance.Port = client.InstanceCreator.FindFreePort()
	instance.Password = uuid.NewRandom().String()

	return client.InstanceCreator.ProvisionInstance(&instance)
}

func (client *SharedNodeClient) UnprovisionInstance(instanceId string) error {
	ok, err := client.InstanceCreator.InstanceExists(instanceId)
	if(err != nil){
		return err
	}
	if !ok {
		return brokerapi.ErrInstanceDoesNotExist
	}
	return client.InstanceCreator.Destroy(instanceId)
}
