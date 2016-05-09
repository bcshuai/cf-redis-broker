package sharednodeagent

import (
	"github.com/pborman/uuid"
	"github.com/pivotal-golang/lager"

	"github.com/bcshuai/brokerapi"
	"github.com/bcshuai/cf-redis-broker/broker"
	"github.com/bcshuai/cf-redis-broker/brokerconfig"
	"github.com/bcshuai/cf-redis-broker/redis"
	sharednode "github.com/bcshuai/cf-redis-broker/sharednodeagentapi"
	"github.com/bcshuai/cf-redis-broker/system"
)

type SharedNodeAgent struct {
	InstanceCreator broker.InstanceCreator
	InstanceRepo    redis.LocalInstanceRepository
	Config          brokerconfig.Config
	Logger          lager.Logger
}

func (client *SharedNodeAgent) Resources() (sharednode.Resource, error) {
	limitation := client.Config.RedisConfiguration.ServiceInstanceLimit
	used, err := client.InstanceRepo.InstanceCount()
	if err != nil {
		return sharednode.Resource{}, err
	}
	instanceStatus := sharednode.ResourceStatus{
		All:  limitation,
		Used: used,
		Free: limitation - used,
	}

	return sharednode.Resource{
		InstanceStatus: instanceStatus,
	}, nil
}

func (client *SharedNodeAgent) AllInstances() ([]*redis.Instance, error) {
	return client.InstanceRepo.AllInstances()
}

func (client *SharedNodeAgent) InstanceInfo(instanceId string) (*redis.Instance, error) {
	return client.InstanceRepo.FindByID(instanceId)
}

func (client *SharedNodeAgent) InstanceExists(instanceId string) (bool, error) {
	return client.InstanceRepo.InstanceExists(instanceId)
}

func (client *SharedNodeAgent) InstanceCredential(instanceId string) (broker.InstanceCredentials, error) {
	ok, err := client.InstanceRepo.InstanceExists(instanceId)
	if err != nil {
		return broker.InstanceCredentials{}, err
	}
	if !ok {
		return broker.InstanceCredentials{}, brokerapi.ErrInstanceDoesNotExist
	}
	instance, err := client.InstanceRepo.FindByID(instanceId)
	if err != nil {
		return broker.InstanceCredentials{}, err
	}

	return broker.InstanceCredentials{
		Host:     instance.Host,
		Port:     instance.Port,
		Password: instance.Password,
	}, nil
}

func (client *SharedNodeAgent) ProvisionInstance(instance redis.Instance) error {
	ok, err := client.InstanceCreator.InstanceExists(instance.ID)
	if err != nil {
		return err
	}
	if ok {
		return brokerapi.ErrInstanceAlreadyExists
	}

	instanceCount, err := client.InstanceRepo.InstanceCount()
	if err != nil {
		return err
	}

	if instanceCount >= client.Config.RedisConfiguration.ServiceInstanceLimit {
		return brokerapi.ErrInstanceLimitMet
	}

	client.Logger.Info("SharedNodeAgent.ProvisonInstance", lager.Data{
		"instance": instance,
	})

	instance.Port, err = system.FindFreePort()
	if err != nil {
		return err
	}
	instance.Password = uuid.NewRandom().String()

	return client.InstanceRepo.ProvisonInstance(&instance)
}

func (client *SharedNodeAgent) UnprovisionInstance(instanceId string) error {
	ok, err := client.InstanceCreator.InstanceExists(instanceId)
	if err != nil {
		return err
	}
	if !ok {
		return brokerapi.ErrInstanceDoesNotExist
	}
	return client.InstanceCreator.Destroy(instanceId)
}
