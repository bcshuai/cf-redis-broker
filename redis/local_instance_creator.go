package redis

import (
	"time"

	"github.com/pborman/uuid"

	"github.com/bcshuai/brokerapi"
	"github.com/bcshuai/cf-redis-broker/brokerconfig"
)

type ProcessController interface {
	StartAndWaitUntilReady(instance *Instance, configPath, instanceDataDir, pidfilePath, logfilePath string, timeout time.Duration) error
	Kill(instance *Instance) error
}

type LocalInstanceRepository interface {
	FindByID(instanceID string) (*Instance, error)
	InstanceExists(instanceID string) (bool, error)
	Setup(instance *Instance) error
	Delete(instanceID string) error
	InstanceDataDir(instanceID string) string
	InstanceConfigPath(instanceID string) string
	InstanceLogFilePath(instanceID string) string
	InstancePidFilePath(instanceID string) string
	InstanceCount() (int, error)
	Lock(instance *Instance) error
	Unlock(instance *Instance) error
}

type LocalInstanceCreator struct {
	LocalInstanceRepository
	FindFreePort       func() (int, error)
	ProcessController  ProcessController
	RedisConfiguration brokerconfig.ServiceConfiguration
}

func (localInstanceCreator *LocalInstanceCreator) Create(instanceID string) error {
	instanceCount, err := localInstanceCreator.InstanceCount()
	if err != nil {
		return err
	}

	if instanceCount >= localInstanceCreator.RedisConfiguration.ServiceInstanceLimit {
		return brokerapi.ErrInstanceLimitMet
	}
    
	port, _ := localInstanceCreator.FindFreePort()
	instance := &Instance{
		ID:       instanceID,
		Port:     port,
		Host:     localInstanceCreator.RedisConfiguration.Host,
		Password: uuid.NewRandom().String(),
	}
	return localInstanceCreator.provisonInstance(instance)
}
func (localInstanceCreator *LocalInstanceCreator) CreateWithRestriction(instanceID string, max_memory_in_mb, max_client_connection int) error {
	instanceCount, err := localInstanceCreator.InstanceCount()
	if err != nil {
		return err
	}

	if instanceCount >= localInstanceCreator.RedisConfiguration.ServiceInstanceLimit {
		return brokerapi.ErrInstanceLimitMet
	}
    
	port, _ := localInstanceCreator.FindFreePort()
	instance := &Instance{
		ID:       instanceID,
		Port:     port,
		Host:     localInstanceCreator.RedisConfiguration.Host,
		Password: uuid.NewRandom().String(),
		MaxClientConnections: max_client_connection,
		MaxMemoryInMB: max_memory_in_mb,
	}
	return localInstanceCreator.provisonInstance(instance)
}

func (localInstanceCreator *LocalInstanceCreator) provisonInstance(instance *Instance) error {

	err := localInstanceCreator.Setup(instance)
	if err != nil {
		return err
	}

	err = localInstanceCreator.startLocalInstance(instance)
	if err != nil {
		return err
	}

	err = localInstanceCreator.Unlock(instance)
	if err != nil {
		return err
	}

	return nil
}

func (localInstanceCreator *LocalInstanceCreator) Destroy(instanceID string) error {
	instance, err := localInstanceCreator.FindByID(instanceID)
	if err != nil {
		return err
	}

	err = localInstanceCreator.Lock(instance)
	if err != nil {
		return err
	}

	err = localInstanceCreator.ProcessController.Kill(instance)
	if err != nil {
		return err
	}

	return localInstanceCreator.Delete(instanceID)
}

func (localInstanceCreator *LocalInstanceCreator) startLocalInstance(instance *Instance) error {
	configPath := localInstanceCreator.InstanceConfigPath(instance.ID)
	instanceDataDir := localInstanceCreator.InstanceDataDir(instance.ID)
	logfilePath := localInstanceCreator.InstanceLogFilePath(instance.ID)
	pidfilePath := localInstanceCreator.InstancePidFilePath(instance.ID)

	timeout := time.Duration(localInstanceCreator.RedisConfiguration.StartRedisTimeoutSeconds) * time.Second
	return localInstanceCreator.ProcessController.StartAndWaitUntilReady(instance, configPath, instanceDataDir, pidfilePath, logfilePath, timeout)
}
