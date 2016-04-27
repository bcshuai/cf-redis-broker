package sharednode

import (
	"sync"
	"sync/atomic"
	"errors"

	"github.com/bcshuai/cf-redis-broker/redis"
	"github.com/bcshuai/cf-redis-broker/broker"
	"github.com/pivotal-golang/lager"
	api "github.com/bcshuai/cf-redis-broker/sharednodeagentapi"
)

type MultiSharedNodeAgentClient struct {
	AgentClients	[]*SharedNodeAgentClient
	Logger	lager.Logger
}
func (client *MultiSharedNodeAgentClient) Resources() (api.Resource, error) {
	nodeCount := len(client.AgentClients)

	var wg sync.WaitGroup
	var lock sync.RWMutex

	wg.Add(nodeCount)

	nodeResourceMap := map[string]api.Resource{}
	nodeErrorMap := map[string]error{}

	for _, c := range client.AgentClients {
		go func(c *SharedNodeAgentClient){
			defer wg.Done()

			res, err := c.Resources()
			lock.Lock()
			nodeResourceMap[c.Host] = res
			nodeErrorMap[c.Host] = err	
			lock.Unlock()
		}(c)
	}
	wg.Wait()

	client.Logger.Info("MultiSharedNodeAgentClient.Resources", lager.Data{
			"resources": nodeResourceMap,
	})

	allResource := api.Resource{
		InstanceStatus: api.ResourceStatus{ 
			All: 0,
			Free: 0,
			Used: 0,
		},
	}
	for host, res := range nodeResourceMap {
		err, _ := nodeErrorMap[host]
		if err != nil {
			continue
		}
		allResource.InstanceStatus.All =  allResource.InstanceStatus.All + res.InstanceStatus.All
		allResource.InstanceStatus.Free = allResource.InstanceStatus.Free + res.InstanceStatus.Free
		allResource.InstanceStatus.Used = allResource.InstanceStatus.Used + res.InstanceStatus.Used
	}
	return allResource, nil
}
func (client *MultiSharedNodeAgentClient) AllInstances() ([]*redis.Instance, error) {
	nodeCount := len(client.AgentClients)

	var wg sync.WaitGroup
	var lock sync.RWMutex

	wg.Add(nodeCount)
	instances := []*redis.Instance{}

	for _, c := range client.AgentClients {
		go func(c *SharedNodeAgentClient){
			defer wg.Done()

			nodeInstances, err := c.AllInstances()
			if(err == nil){
				
				for _, i := range nodeInstances {
					lock.Lock()
					instances = append(instances, i)
					lock.Unlock()
				}
				
			}
		}(c)
	}
	wg.Wait()
	return instances, nil
}
func (client *MultiSharedNodeAgentClient) InstanceInfo(instanceId string) (*redis.Instance, error){
	nodeCount := len(client.AgentClients)

	var wg sync.WaitGroup
	var lock sync.RWMutex

	wg.Add(nodeCount)
	instances := []*redis.Instance{}
	for _, c := range client.AgentClients {
		go func(c *SharedNodeAgentClient){
			defer wg.Done()

			nodeInstance, err := c.InstanceInfo(instanceId)
			if(err == nil){
				lock.Lock()
				instances = append(instances, nodeInstance)
				lock.Unlock()
			}
		}(c)
	}
	wg.Wait()

	count := len(instances)
	if(count == 0){
		return nil, nil
	}
	if(count > 1) {
		return instances[0], errors.New("too many instances")
	}

	return instances[0], nil
}
func (client *MultiSharedNodeAgentClient) InstanceExists(instanceId string) (bool, error) {
	nodeCount := len(client.AgentClients)

	var wg sync.WaitGroup

	wg.Add(nodeCount)

	var count uint32 = 0

	for _, c := range client.AgentClients {
		go func(tc *SharedNodeAgentClient) {
			defer wg.Done()
            if(tc == nil){
            	client.Logger.Info("InstanceExists doing", lager.Data{
					"Msg": "why client is nil",
				})
            }
			ok, err := tc.InstanceExists(instanceId)
			if (err == nil && ok){
				atomic.AddUint32(&count, 1)
			}
		}(c)
	}
	wg.Wait()
	
	client.Logger.Info("InstanceExists done", lager.Data{
			"instanceId": instanceId,
	})
	if(count > 1){
		return true, errors.New("too many instances")
	} else if (count == 0){
		client.Logger.Info("Instance does not exist", lager.Data{
			"instanceId": instanceId,
		})
		return false, nil
	} else {
		return true, nil
	}
}
func (client *MultiSharedNodeAgentClient) InstanceCredential(instanceId string) (broker.InstanceCredentials, error) {
	nodeCount := len(client.AgentClients)

	var wg sync.WaitGroup
	var lock sync.RWMutex

	wg.Add(nodeCount)
	credentials := []broker.InstanceCredentials{}

	for _, c := range client.AgentClients {
		go func(c *SharedNodeAgentClient){
			defer wg.Done()

			nodeCredential, err := c.InstanceCredential(instanceId)
			if(err == nil){
				lock.Lock()
				credentials = append(credentials, nodeCredential)
				lock.Unlock()
			}
		}(c)
	}
	wg.Wait()

	count := len(credentials)
	if(count == 0){
		client.Logger.Info("Instance does not exist", lager.Data{
			"instanceId": instanceId,
		})
		return broker.InstanceCredentials{}, errors.New("instance does not exist")
	}
	if(count > 1) {
		return broker.InstanceCredentials{}, errors.New("too many instances")
	}

	return credentials[0], nil
}
func (client *MultiSharedNodeAgentClient) getClinetByHost(host string) *SharedNodeAgentClient {
	for _, c := range client.AgentClients {
		if(c.Host == host){
			return c
		}
	}
	return nil
}
func (client *MultiSharedNodeAgentClient) ProvisionInstance(instance redis.Instance) error {
	nodeCount := len(client.AgentClients)

	var wg sync.WaitGroup
	var lock sync.RWMutex

	wg.Add(nodeCount)

	nodeResourceMap := map[string]api.Resource{}
	for _, c := range client.AgentClients {
		go func(c *SharedNodeAgentClient){
			defer wg.Done()

			res, err := c.Resources()
			if(err == nil){
				lock.Lock()
				client.Logger.Info("resource: ", lager.Data{
					"resources": res,
				})
				nodeResourceMap[c.Host] = res
				lock.Unlock()
			} else {
				client.Logger.Info("get resource failed: ", lager.Data{
					"err": err,
				})
			}
		}(c)
	}
	wg.Wait()

	client.Logger.Info("MultiSharedNodeAgentClient.Resources", lager.Data{
			"resources": nodeResourceMap,
	})

	targetHost := ""
	var minUsageRate float32 = -1.0
	for host, res := range nodeResourceMap {
		if res.InstanceStatus.All == 0 {
			continue
		}

		usage := float32(res.InstanceStatus.Free) / float32(res.InstanceStatus.All)
		client.Logger.Info("Node.Resource", lager.Data{
			host: res,
			"free": usage,
		})
		if(usage > minUsageRate){
			minUsageRate = usage
			client.Logger.Info("Node.Resource", lager.Data{
				"minUsageRate": minUsageRate,
			})

			targetHost = host
			client.Logger.Info("Node.Resource", lager.Data{
				"target Host": targetHost,
			})
		}
	}
	if(targetHost == ""){
		return errors.New("unable to find available node to provision the instance")
	}
	targetClient := client.getClinetByHost(targetHost)
	return targetClient.ProvisionInstance(instance)
}
func (client *MultiSharedNodeAgentClient) UnprovisionInstance(instanceId string) error {
	instance, err := client.InstanceInfo(instanceId)
	if(err != nil){
		return err
	}
	if(instance == nil){
		return nil
	}

	c := client.getClinetByHost(instance.Host)
	return c.UnprovisionInstance(instanceId)
}
