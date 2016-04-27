package sharednodeagentapi

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/gorilla/mux"
	"github.com/bcshuai/cf-redis-broker/redis"
	"github.com/bcshuai/cf-redis-broker/broker"
)
type ResourceStatus struct {
	All    	int		`json:"all"`
	Used 	int		`json:"used"`
	Free 	int 		`json:"free"`
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

func New(provider ApiProvider) http.Handler {
	router := mux.NewRouter()

	router.Path("/resources").
			Methods("GET").
			HandlerFunc(resourcesHandler(provider))
	router.Path("/all_instances").
			Methods("GET").
			HandlerFunc(allInstancesHandler(provider))
	router.Path("/instance/{id}").
			Methods("GET").
			HandlerFunc(instanceInfoHandler(provider))
	router.Path("/exist/{id}").
			Methods("GET").
			HandlerFunc(instanceExistsHandler(provider))
	router.Path("/credential/{id}").
			Methods("GET").
			HandlerFunc(instanceCredentialHandler(provider))
	router.Path("/instance/{id}").
			Methods("PUT").
			HandlerFunc(provisionInstanceHandler(provider))
	router.Path("/instance/{id}").
			Methods("DELETE").
			HandlerFunc(unprovisionInstanceHandler(provider))
	return router
}

func resourcesHandler(provider ApiProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := provider.Resources()
		if(err != nil){
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		encoder.Encode(res)
	}
}

func allInstancesHandler(provider ApiProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		instances, err := provider.AllInstances()
		if(err != nil){
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		encoder.Encode(instances)
	}
}

func instanceInfoHandler(provider ApiProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		instanceId, ok := vars["id"]; 
		if !ok {
			http.Error(w, "Instance ID is required", http.StatusInternalServerError)
			return
		}
		instance, err := provider.InstanceInfo(instanceId)
		if(err != nil){
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		encoder.Encode(instance)
	}	
}

func instanceExistsHandler(provider ApiProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		instanceId, ok := vars["id"]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		ok, err := provider.InstanceExists(instanceId)
		if !ok || err != nil {
			w.WriteHeader(http.StatusNotFound)
		}else{
			w.WriteHeader(http.StatusOK)
		}
	}
}

func instanceCredentialHandler(provider ApiProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		instanceId, ok := vars["id"]
		if !ok {
			http.Error(w, "Instance ID is required", http.StatusInternalServerError)
			return
		}
		credential, err := provider.InstanceCredential(instanceId)
		if(err != nil){
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		encoder.Encode(credential)
	}
}

func provisionInstanceHandler(provider ApiProvider) http.HandlerFunc {
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

		log.Println(string(body))
		err = provider.ProvisionInstance(instance)
		if(err != nil){
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

func unprovisionInstanceHandler(provider ApiProvider) http.HandlerFunc {
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

