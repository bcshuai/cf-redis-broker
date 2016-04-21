package agentapi-sharednode

import (
	"net/http"
	"encoding/json"
	"io/ioutil"

	"github.com/gorilla/mux"
	"github.com/bcshuai/cf-redis-broker/redis"
	"github.com/bcshuai/cf-redis-broker/broker"
)
type ResourceStatus struct {
	All    	uint64		`json:"all"`
	Used 	uint64		`json:"used"`
	Free 	uint64 		`json:"free"`
}

type Resource struct {
	//MemoryStatus 	ResourceStatus 		`json:"memory"`
	//DiskStatus   	ResourceStatus 		`json:"disk"`
	InstanceStatus 	ResourceStatus 		`json:"instances"`
	//CPUS			int 				`json:"cpus"`
}

type ApiProvider interface {
	Resources() (Resource, error)
	AllInstances() ([]*redis.Instance, error)
	InstanceInfo(instanceId string) (*redis.Instance, error)
	InstanceExists(instanceId string) (bool, error)
	InstanceCredential(instanceId string) (broker.InstanceCredentials, error)
	ProvisionInstance(instance redis.Instance) error
	UnprovisionInstance(instanceId string) error
}

func New(provider *ApiProvider) http.Handler {
	router := mux.NewRouter()

	router.Path("/resources").
			Methods("GET").
			HandlerFunc(ResourcesHandler(provider))
	router.Path("/all_instances").
			Methods("GET").
			HandlerFunc(AllInstancesHandler(provider))
	router.Path("/instance/{id}").
			Methods("GET").
			HandlerFunc(InstanceInfoHandler(provider))
	router.Path("/instance/{id}").
			Methods("HEAD").
			HandlerFunc(InstanceExistsHandler(provider))
	router.Path("/credential/{id}").
			Methods("GET").
			HandlerFunc(InstanceCredentialHandler(provider))
	router.Path("/instance/{id}").
			Methods("PUT").
			HandlerFunc(ProvisionInstanceHandler(provider))
	router.Path("/instance/{id}").
			Methods("DELETE").
			HandlerFunc(UnprovisionInstanceHandler(provider))
	return router
}

func ResourcesHandler(provider *ApiProvider) http.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := provider.Resources()
		if(err != nil){
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		encoder.Encode(credentials)
	}
}

func AllInstancesHandler(provider *ApiProvider) http.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		instances, err := provider.AllInstances()
		if(err != nil){
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		encoder.Encode(instances)
	}
}

func InstanceInfoHandler(provider *ApiProvider) http.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		instanceId, ok := vars["id"]; 
		if !ok {
			http.Error(w, "Instance ID is required", http.StatusInternalServerError)
			return
		}
		instance, err := provider.InstanceInfo(instanceId)
		w.Header.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		encoder.Encode(instance)
	}	
}

func InstanceExistsHandler(provider *ApiProvider) http.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if instanceId, ok := vars["id"]; !ok {
			http.Error(w, "Instance ID is required", http.StatusInternalServerError)
			return
		}
		ok, err := provider.InstanceExists(instanceId)
		if(err != nil){
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if !ok {
			w.WriteHeader(http.StatusNotFound)
		}else{
			w.WriteHeader(http.StatusOK)
		}
	}
}

func InstanceCredentialHandler(provider *ApiProvider) http.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if instanceId, ok := vars["id"]; !ok {
			http.Error(w, "Instance ID is required", http.StatusInternalServerError)
			return
		}
		credential, err := provider.InstanceCredential(instanceId)
		if(err != nil){
			http.Error(w, err.Error, http.StatusInternalServerError)
			return
		}
		w.Header.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		encoder.Encode(credential)
	}
}

func ProvisionInstanceHandler(provider *ApiProvider) http.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if(err != nil){
			http.Error(w, "Server Internal Error", http.StatusInternalServerError)
			return
		}
		instance := redis.Instance{}
		err = json.Unmarshal(body, &instance)
		if err != nil {
			http.Error(w, "The request contains wrong format content", http.StatusInternalServerError)
			return
		}
		err := provider.ProvisionInstance(instance)
		if(err != nil){
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

func UnprovisionInstanceHandler(provider *ApiProvider) http.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		instanceId, ok := vars["id"]
		if !ok {
			http.Error(w, "Instance ID is required", http.StatusInternalServerError)
			return
		}
		err := provider.UnprovisionInstance(instanceId)
		if(err != nil){
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

