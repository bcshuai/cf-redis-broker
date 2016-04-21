package redis

import (
	"github.com/bcshuai/cf-redis-broker/brokerconfig"
)

type SharedRemoteRepository struct {
	Config   	brokerconfig.Config
}

func NewShareRemoteRepository(config brokerconfig.Config) (*SharedRemoteRepository, error) {
	
}

/* implements interface InstanceCreator */
func (repo *SharedRemoteRepository) Create(instanceID string) error {

}
func (repo *SharedRemoteRepository) CreateWithRestriction(instanceID string, max_mem_in_mb int, max_client_connection int) error {

}
func (repo *SharedRemoteRepository) Destroy(instanceID string) error {

}
func (repo *SharedRemoteRepository) InstanceExists(instanceID string) (bool, error) {

}

/* implements interface InstanceBinder */
func (repo *SharedRemoteRepository) Bind(instanceID string, bindingID string) (InstanceCredentials, error) {

}

func (repo *SharedRemoteRepository) Unbind(instanceID string, bindingID string) error {

}

func (repo *SharedRemoteRepository) InstanceExists(instanceID string) (bool, error) {

}

/* other method needed by this framework */
func (repo *SharedRemoteRepository) FindByID(instanceID string) (*Instance, error){

}
func (repo *SharedRemoteRepository) AllInstances() ([]*Instance, error) {

}
func (repo *SharedRemoteRepository) InstanceLimit() int {

}
func (repo *SharedRemoteRepository) AvailableInstances() []*Instance {

}
func (repo *SharedRemoteRepository) BindingsForInstance(instanceID string) ([]string, error){

}
