package sharednode

import (
	"strconv"

	"github.com/bcshuai/cf-redis-broker/broker"
	"github.com/bcshuai/cf-redis-broker/brokerconfig"
	"github.com/bcshuai/cf-redis-broker/redis"
	api "github.com/bcshuai/cf-redis-broker/sharednodeagentapi"
	"github.com/pivotal-golang/lager"
)

type SharedRemoteRepository struct {
	AgentClient api.ApiProvider
	Logger      lager.Logger
}

func NewShareRemoteRepository(config brokerconfig.Config, logger lager.Logger) (*SharedRemoteRepository, error) {

	agentPort, err := strconv.Atoi(config.AgentPort)
	if err != nil {
		return nil, err
	}

	agentCredential := config.AuthConfiguration

	clients := []*SharedNodeAgentClient{}

	//for _, host := range config.RedisConfiguration.Shared.Nodes {
	for _, host := range config.Shared.Nodes {
		client := &SharedNodeAgentClient{
			Host:            host,
			Port:            agentPort,
			AgentCredential: agentCredential,
		}
		clients = append(clients, client)
	}

	multiAgentClient := NewMultiSharedNodeAgentClient(clients, logger)

	return &SharedRemoteRepository{
		AgentClient: multiAgentClient,
		Logger:      logger,
	}, nil
}

/* implements interface InstanceCreator */
func (repo *SharedRemoteRepository) Create(instanceID string) error {
	//no op
	return nil
}

func (repo *SharedRemoteRepository) CreateWithRestriction(instanceID string, max_mem_in_mb int, max_client_connection int) error {
	instance := redis.Instance{
		ID:                   instanceID,
		MaxMemoryInMB:        max_mem_in_mb,
		MaxClientConnections: max_client_connection,
	}
	repo.Logger.Info("SharedRemoteRepository.CreateWithRestriction", lager.Data{
		"instance": instance,
	})
	return repo.AgentClient.ProvisionInstance(instance)
}
func (repo *SharedRemoteRepository) Destroy(instanceID string) error {
	return repo.AgentClient.UnprovisionInstance(instanceID)
}
func (repo *SharedRemoteRepository) InstanceExists(instanceID string) (bool, error) {
	return repo.AgentClient.InstanceExists(instanceID)
}

/* implements interface InstanceBinder */
func (repo *SharedRemoteRepository) Bind(instanceID string, bindingID string) (broker.InstanceCredentials, error) {
	return repo.AgentClient.InstanceCredential(instanceID)
}

func (repo *SharedRemoteRepository) Unbind(instanceID string, bindingID string) error {
	//no op
	return nil
}

/* other method needed by this framework */
func (repo *SharedRemoteRepository) FindByID(instanceID string) (*redis.Instance, error) {
	return repo.AgentClient.InstanceInfo(instanceID)
}
func (repo *SharedRemoteRepository) AllInstances() ([]*redis.Instance, error) {
	return repo.AgentClient.AllInstances()
}
func (repo *SharedRemoteRepository) InstanceLimit() int {
	res, err := repo.AgentClient.Resources()
	if err != nil {
		return 0
	}
	return res.InstanceStatus.All
}
func (repo *SharedRemoteRepository) AvailableInstances() []*redis.Instance {
	return nil
}
func (repo *SharedRemoteRepository) BindingsForInstance(instanceID string) ([]string, error) {
	return []string{}, nil
}
