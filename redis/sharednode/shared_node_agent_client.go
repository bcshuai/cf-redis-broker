package sharednode

import (
	"fmt"
	"errors"
	"crypto/tls"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"bufio"
	"strings"
	"strconv"

	"github.com/bcshuai/cf-redis-broker/brokerconfig"
	"github.com/bcshuai/cf-redis-broker/broker"
	"github.com/bcshuai/cf-redis-broker/redis"
	sharednode "github.com/bcshuai/cf-redis-broker/sharednodeagentapi"
)

type SharedNodeAgentClient struct {
	Host				string
	Port				int
	AgentCredential 	brokerconfig.AuthConfiguration
}

func (client *SharedNodeAgentClient) getEndpoint(path string) string {
	portStr := strconv.Itoa(client.Port)
	return "https://" + client.Host + ":" + portStr + path
}
func formatErrResponse(response *http.Response) error {
	//body, _ := ioutil.ReadAll(response.Body)
	formattedBody := ""
	// if len(body) > 0 {
	// 	formattedBody = fmt.Sprintf(", %s", string(body))
	// }
	return errors.New(fmt.Sprintf("Agent error: %d%s", response.StatusCode, formattedBody))
}
func (client *SharedNodeAgentClient) Resources() (sharednode.Resource, error){
	const (
		PATH = "/resources"
		METHOD = "GET"
	)
	resource := sharednode.Resource{}
	response, err := client.doAuthenticatedRequest(PATH, METHOD, nil)
	if(err != nil){
		return resource, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if(err != nil) {
		return resource, err
	}
	if response.StatusCode != http.StatusOK {
		return resource, formatErrResponse(response)
	}

	err = json.Unmarshal(body, &resource)
	return resource, err
}
func (client *SharedNodeAgentClient) AllInstances() ([]*redis.Instance, error){
	const (
		PATH = "/all_instances"
		METHOD = "GET"
	)
	instances := []*redis.Instance{}
	response, err := client.doAuthenticatedRequest(PATH, METHOD, nil)
	if(err != nil) {
		return instances, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if(err != nil){
		return instances, err
	}

	if(response.StatusCode != http.StatusOK){
		return instances, formatErrResponse(response)
	}

	err = json.Unmarshal(body, &instances)
	if(err != nil){
		return instances, err
	}
	for _, instance := range instances {
		instance.Host = client.Host
	}
	return instances, nil
}
func (client *SharedNodeAgentClient) InstanceInfo(instanceId string) (*redis.Instance, error){
	const (
		PATH = "/instance/"
		METHOD = "GET"
	)

	var instance *redis.Instance = nil
	response, err := client.doAuthenticatedRequest(PATH + instanceId, METHOD, nil)
	if(err != nil) {
		return instance, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if(err != nil){
		return instance, err
	}

	if(response.StatusCode != http.StatusOK){
		return instance, formatErrResponse(response)
	}

	err = json.Unmarshal(body, &instance)
	if(err != nil){
		return &redis.Instance{}, err
	}
	instance.Host = client.Host
	return instance, nil
}
func (client *SharedNodeAgentClient) InstanceExists(instanceId string) (bool, error) {
	const (
		PATH = "/exist/"
		METHOD = "GET"
	)

	response, err := client.doAuthenticatedRequest(PATH + instanceId, METHOD, nil)
	if(err != nil) {
		return false, err
	}
	defer response.Body.Close()

	if(response.StatusCode == http.StatusNotFound){
		return false, nil
	} else if (response.StatusCode == http.StatusOK){
		return true, nil
	} else {
		return false, formatErrResponse(response)
	}
}

func (client *SharedNodeAgentClient) InstanceCredential(instanceId string) (broker.InstanceCredentials, error) {
	const (
		PATH = "/credential/"
		METHOD = "GET"
	)

	credential := broker.InstanceCredentials{}
	response, err := client.doAuthenticatedRequest(PATH + instanceId, METHOD, nil)
	if(err != nil){
		return credential, err
	}
	defer response.Body.Close()

	if(response.StatusCode != http.StatusOK){
		return credential, formatErrResponse(response)
	}

	body, err := ioutil.ReadAll(response.Body)
	if(err != nil){
		return credential, err
	}

	err = json.Unmarshal(body, &credential)
	if( err != nil ){
		return credential, err
	} else {
		credential.Host = client.Host
		return credential, nil
	}
}
func (client *SharedNodeAgentClient) ProvisionInstance(instance redis.Instance) error {
	const (
		PATH = "/instance/"
		METHOD = "PUT"
	)

	instanceData, err := json.Marshal(instance)
	if(err != nil){
		return err
	}
	strReader := strings.NewReader(string(instanceData))

	strBufReader := bufio.NewReader(strReader)

	response, err := client.doAuthenticatedRequest(PATH + instance.ID, METHOD, strBufReader)
	if(err != nil){
		return err
	}
	defer response.Body.Close()

	if(response.StatusCode != http.StatusOK){
		return formatErrResponse(response)
	} else {
		return nil
	}
}
func (client *SharedNodeAgentClient) UnprovisionInstance(instanceId string) error {
	const (
		PATH = "/instance/"
		METHOD = "DELETE"
	)

	response, err := client.doAuthenticatedRequest(PATH + instanceId, METHOD, nil)
	if(err != nil){
		return err
	}
	defer response.Body.Close()

	if(response.StatusCode != http.StatusOK){
		return formatErrResponse(response)
	} else {
		return nil
	}
}

func (client *SharedNodeAgentClient) doAuthenticatedRequest(path, method string, body *bufio.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method, client.getEndpoint(path), body)
	if(err != nil){
		return nil, err
	}
	request.SetBasicAuth(client.AgentCredential.Username, client.AgentCredential.Password)

	httpClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	return httpClient.Do(request)
}