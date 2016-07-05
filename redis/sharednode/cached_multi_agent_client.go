package sharednode

import (
	"errors"

	"github.com/bcshuai/cf-redis-broker/broker"
	"github.com/bcshuai/cf-redis-broker/redis"
	"github.com/bcshuai/cf-redis-broker/redis/sharednode/cache"
	"github.com/pivotal-golang/lager"

	api "github.com/bcshuai/cf-redis-broker/sharednodeagentapi"
)

func getClientByHost(clients []*SharedNodeAgentClient, host string) *SharedNodeAgentClient {
	for _, c := range clients {
		if c.Host == host {
			return c
		}
	}
	return nil
}

type CachedMultiSharedNodeAgentClient struct {
	clients        []*SharedNodeAgentClient
	resourceCache  *cache.Cache
	instancesCache *cache.Cache
	logger         lager.Logger
}

type resourceCacheUpdater struct {
	Clients []*SharedNodeAgentClient
}
type instancesCacheUpdater struct {
	Clients []*SharedNodeAgentClient
}

func (updater *resourceCacheUpdater) Keys() []string {
	var hosts []string = []string{}
	for _, client := range updater.Clients {
		hosts = append(hosts, client.Host)
	}
	return hosts
}
func (updater *resourceCacheUpdater) Value(key string) interface{} {
	nodeClient := getClientByHost(updater.Clients, key)
	if nodeClient == nil {
		return nil
	}
	res, err := nodeClient.Resources()
	if err != nil {
		return nil
	} else {
		return res
	}
}

func (updater *instancesCacheUpdater) Keys() []string {
	var hosts []string = []string{}
	for _, client := range updater.Clients {
		hosts = append(hosts, client.Host)
	}
	return hosts
}
func (updater *instancesCacheUpdater) Value(key string) interface{} {
	nodeClient := getClientByHost(updater.Clients, key)
	if nodeClient == nil {
		return nil
	}
	instances, err := nodeClient.AllInstances()
	if err != nil {
		return nil
	} else {
		return instances
	}
}

func NewCachedMultiSharedNodeAgentClient(clients []*SharedNodeAgentClient, logger lager.Logger) *CachedMultiSharedNodeAgentClient {

	resourceCache := cache.NewCache(true, 50, &resourceCacheUpdater{
		Clients: clients,
	})
	instancesCache := cache.NewCache(true, 50, &instancesCacheUpdater{
		Clients: clients,
	})

	cachedMultiSharedClient := &CachedMultiSharedNodeAgentClient{
		clients:        clients,
		resourceCache:  resourceCache,
		instancesCache: instancesCache,
		logger:         logger,
	}

	resourceCache.Start()
	instancesCache.Start()

	return cachedMultiSharedClient
}

func (client *CachedMultiSharedNodeAgentClient) Resources() (api.Resource, error) {
	allResource := api.Resource{
		InstanceStatus: api.ResourceStatus{
			All:  0,
			Free: 0,
			Used: 0,
		},
		MemoryStatus: api.ResourceStatus{
			All:  0,
			Free: 0,
			Used: 0,
		},
	}
	for _, node := range client.clients {
		host := node.Host
		val := client.resourceCache.Get(host)
		if val == nil {
			client.logger.Info("CachedMultisharedClient.Resources", lager.Data{
				"Msg": "failed to get resource info from host: " + host,
			})
			continue
		}
		res, ok := val.(api.Resource)
		if !ok {
			client.logger.Info("CachedMultisharedClient.Resources", lager.Data{
				"Msg":   "unable to convert interface{} to api.Resource",
				"Value": val,
			})
			continue
		}
		allResource.InstanceStatus.All += res.InstanceStatus.All
		allResource.InstanceStatus.Free += res.InstanceStatus.Free
		allResource.InstanceStatus.Used += res.InstanceStatus.Used

		allResource.MemoryStatus.All += res.MemoryStatus.All
		allResource.MemoryStatus.Free += res.MemoryStatus.Free
		allResource.MemoryStatus.Used += res.MemoryStatus.Used
	}
	return allResource, nil
}
func (client *CachedMultiSharedNodeAgentClient) AllInstances() ([]*redis.Instance, error) {
	allInstances := []*redis.Instance{}

	for _, node := range client.clients {
		val := client.instancesCache.Get(node.Host)
		if val == nil {
			client.logger.Info("CachedMultisharedClient.AllInstances", lager.Data{
				"Msg": "failed to get All instances from host: " + node.Host,
			})
			continue
		}
		nodeInstances, ok := val.([]*redis.Instance)
		if !ok {
			client.logger.Info("CachedMultisharedClient.AllInstances", lager.Data{
				"Msg":   "unable to convert interface{} to []*redis.Instance",
				"Value": val,
			})
			continue
		}
		for _, instance := range nodeInstances {
			allInstances = append(allInstances, instance)
		}
	}

	return allInstances, nil
}

func (client *CachedMultiSharedNodeAgentClient) InstanceInfo(instanceId string) (*redis.Instance, error) {
	candi := []*redis.Instance{}
	allInstances, _ := client.AllInstances()
	for _, ins := range allInstances {
		if ins.ID == instanceId {
			candi = append(candi, ins)
		}
	}
	count := len(candi)
	if count == 0 {
		return nil, nil
	}
	if count > 1 {
		return candi[0], errors.New("too many instances")
	}

	return candi[0], nil
}

func (client *CachedMultiSharedNodeAgentClient) InstanceExists(instanceId string) (bool, error) {
	candi := []*redis.Instance{}
	allInstances, _ := client.AllInstances()
	for _, ins := range allInstances {
		if ins.ID == instanceId {
			candi = append(candi, ins)
		}
	}
	count := len(candi)
	if count == 0 {
		return false, nil
	}
	if count > 1 {
		return true, errors.New("too many instances")
	}

	return true, nil
}

func (client *CachedMultiSharedNodeAgentClient) InstanceCredential(instanceId string) (broker.InstanceCredentials, error) {
	candi := []*redis.Instance{}
	allInstances, _ := client.AllInstances()
	for _, ins := range allInstances {
		if ins.ID == instanceId {
			candi = append(candi, ins)
		}
	}
	count := len(candi)
	if count == 0 {
		return broker.InstanceCredentials{}, errors.New("instance does not exist")
	}
	if count > 1 {
		return broker.InstanceCredentials{}, errors.New("too many instances")
	}

	return broker.InstanceCredentials{
		Host:     candi[0].Host,
		Port:     candi[0].Port,
		Password: candi[0].Password,
	}, nil

}

func (client *CachedMultiSharedNodeAgentClient) ProvisionInstance(instance redis.Instance) error {
	var minUsageRate float32 = -1.0
	var targetNode *SharedNodeAgentClient = nil
	for _, node := range client.clients {
		host := node.Host
		val := client.resourceCache.Get(host)
		if val == nil {
			continue
		}
		res, ok := val.(api.Resource)
		if !ok {
			continue
		}
		if res.InstanceStatus.All != 0 && res.InstanceStatus.Free == 0 {
			continue
		}
		if res.MemoryStatus.Free < instance.MaxMemoryInMB {
			continue
		}
		usage := float32(res.MemoryStatus.Free) / float32(res.MemoryStatus.All)
		client.logger.Info("Node.Resource", lager.Data{
			host:   res,
			"free": usage,
		})
		if usage > minUsageRate {
			minUsageRate = usage
			targetNode = node
		}
	}
	if targetNode == nil {
		return errors.New("unable to find available node to provision the instance")
	}
	return targetNode.ProvisionInstance(instance)
}

func (client *CachedMultiSharedNodeAgentClient) UnprovisionInstance(instanceId string) error {
	instance, err := client.InstanceInfo(instanceId)
	if err != nil {
		return err
	}
	if instance == nil {
		return nil
	}

	node := getClientByHost(client.clients, instance.Host)
	return node.UnprovisionInstance(instanceId)
}
