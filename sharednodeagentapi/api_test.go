package sharednodeagentapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/bcshuai/cf-redis-broker/broker"
	"github.com/bcshuai/cf-redis-broker/redis"
	"github.com/pivotal-golang/lager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func fakeRedisInstance() *redis.Instance {
	fakeInstance := &redis.Instance{
		ID:                   "123456789",
		Host:                 "1.1.1.1",
		Port:                 8888,
		Password:             "password",
		MaxMemoryInMB:        1024,
		MaxClientConnections: 10,
	}
	return fakeInstance
}

type fakeApiProvider struct {
}

func (fake *fakeApiProvider) Resources() (Resource, error) {
	return Resource{
		InstanceStatus: ResourceStatus{
			All:  1,
			Used: 0,
			Free: 1,
		},
	}, nil
}

func (fake *fakeApiProvider) AllInstances() ([]*redis.Instance, error) {
	fakeInstance := fakeRedisInstance()
	instances := []*redis.Instance{fakeInstance}

	return instances, nil
}

func (fake *fakeApiProvider) InstanceInfo(instanceId string) (*redis.Instance, error) {
	fakeInstance := fakeRedisInstance()
	if fakeInstance.ID != instanceId {
		return nil, errors.New("Not Found")
	}
	return fakeInstance, nil
}

func (fake *fakeApiProvider) InstanceExists(instanceId string) (bool, error) {
	fakeInstance := fakeRedisInstance()

	if fakeInstance.ID != instanceId {
		return false, nil
	}
	return true, nil
}

func (fake *fakeApiProvider) InstanceCredential(instanceId string) (broker.InstanceCredentials, error) {
	fakeInstance := fakeRedisInstance()
	if fakeInstance.ID != instanceId {
		return broker.InstanceCredentials{}, errors.New("Not Found")
	}
	return broker.InstanceCredentials{
		Host:     fakeInstance.Host,
		Port:     fakeInstance.Port,
		Password: fakeInstance.Password,
	}, nil
}

func (fake *fakeApiProvider) ProvisionInstance(instance redis.Instance) error {
	if instance.ID == "123456789" {
		return errors.New("Already Exists")
	}
	return nil
}

func (fake *fakeApiProvider) UnprovisionInstance(instanceId string) error {
	fakeInstance := fakeRedisInstance()
	if fakeInstance.ID != instanceId {
		return errors.New("Not Found")
	}
	return nil
}

func makeRequest(method string, url string, body io.Reader) *http.Response {
	request, err := http.NewRequest(method, url, body)
	Ω(err).ShouldNot(HaveOccurred())

	response, err := http.DefaultClient.Do(request)
	Ω(err).ShouldNot(HaveOccurred())

	return response
}

var _ = Describe("redis shared node agent HTTP API", func() {
	var server *httptest.Server
	var apiProvider *fakeApiProvider
	var response *http.Response
	var logger lager.Logger
	BeforeEach(func() {

	})

	JustBeforeEach(func() {
		logger = lager.NewLogger("shared-node-redis-agent")
		logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
		logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))

		apiProvider = &fakeApiProvider{}
		handler := New(apiProvider, logger)
		server = httptest.NewServer(handler)
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Get /resources", func() {
		JustBeforeEach(func() {
			response = makeRequest("GET", server.URL+"/resources", nil)
		})
		Context("When we load the resource info successfully", func() {
			It("return the correct resource info", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).ShouldNot(HaveOccurred())
				var res Resource
				err = json.Unmarshal(body, &res)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(res.InstanceStatus.All).Should(Equal(int(1)))
				Expect(res.InstanceStatus.Free).Should(Equal(int(1)))
				Expect(res.InstanceStatus.Used).Should(Equal(int(0)))
			})
		})
	})
	Describe("Get /all_instances", func() {
		JustBeforeEach(func() {
			response = makeRequest("GET", server.URL+"/all_instances", nil)
		})
		Context("When we load all the instances successfully", func() {
			It("return all the instances", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).ShouldNot(HaveOccurred())

				var instances []*redis.Instance
				err = json.Unmarshal(body, &instances)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(len(instances)).Should(Equal(1))
				Expect(instances[0].ID).Should(Equal("123456789"))
			})
		})
	})
	Describe("Get target instance /instance/{id}", func() {
		Context("When the target instance exists", func() {
			JustBeforeEach(func() {
				response = makeRequest("GET", server.URL+"/instance/123456789", nil)
			})
			It("return the target instance", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).ShouldNot(HaveOccurred())

				var instance redis.Instance
				err = json.Unmarshal(body, &instance)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(instance.ID).Should(Equal("123456789"))
			})
		})
		Context("When the target does not exist", func() {
			JustBeforeEach(func() {
				response = makeRequest("GET", server.URL+"/instance/12345", nil)
			})
			It("return 500", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(response.StatusCode).Should(Equal(500))
				Expect(string(body)).Should(Equal("Not Found\n"))
			})
		})
	})
	Describe("Check the existence of an instance /exist/{id}", func() {
		Context("When the target exists", func() {
			JustBeforeEach(func() {
				response = makeRequest("GET", server.URL+"/exist/123456789", nil)
			})
			It("return true", func() {
				Expect(response.StatusCode).Should(Equal(200))
			})
		})

		Context("When the target does not exist", func() {
			JustBeforeEach(func() {
				response = makeRequest("GET", server.URL+"/exist/123456", nil)
			})
			It("return false", func() {
				Expect(response.StatusCode).Should(Equal(404))
			})
		})
	})
	Describe("Get credential info of an instance /credential/{id}", func() {
		Context("The target instance not exist", func() {
			JustBeforeEach(func() {
				response = makeRequest("GET", server.URL+"/credential/123456", nil)
			})
			It("return 500", func() {
				Expect(response.StatusCode).Should(Equal(500))
			})
		})
		Context("The target instance exists", func() {
			JustBeforeEach(func() {
				response = makeRequest("GET", server.URL+"/credential/123456789", nil)
			})
			It("return the right credential", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).ShouldNot(HaveOccurred())

				var credential broker.InstanceCredentials
				err = json.Unmarshal(body, &credential)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(credential.Host).Should(Equal("1.1.1.1"))
				Expect(credential.Port).Should(Equal(8888))
				Expect(credential.Password).Should(Equal("password"))
			})
		})
	})
	Describe("Provision an instance Put /instance/{id}", func() {
		Context("the provisioning instance already exists", func() {
			JustBeforeEach(func() {
				fakeInstance := fakeRedisInstance()
				fakeInstance.ID = "123456789"
				instanceData, err := json.Marshal(fakeInstance)
				Expect(err).ShouldNot(HaveOccurred())
				response = makeRequest("PUT", server.URL+"/instance/123456789", bytes.NewBuffer(instanceData))
			})
			It("return 500", func() {
				Expect(response.StatusCode).Should(Equal(500))
			})
		})
		Context("the provisioning instance not exist", func() {
			JustBeforeEach(func() {
				fakeInstance := fakeRedisInstance()
				fakeInstance.ID = "123456"
				instanceData, err := json.Marshal(fakeInstance)
				Expect(err).ShouldNot(HaveOccurred())

				response = makeRequest("PUT", server.URL+"/instance/123456", bytes.NewBuffer(instanceData))
			})
			It("the instance is provisioned", func() {
				Expect(response.StatusCode).Should(Equal(200))
			})
		})
	})
	Describe("Unprovision an instance DELETE /instance/{id}", func() {
		Context("the target instance not exist", func() {
			JustBeforeEach(func() {
				response = makeRequest("DELETE", server.URL+"/instance/123456", nil)
			})
			It("return 500", func() {
				Expect(response.StatusCode).Should(Equal(500))
			})
		})
		Context("the target instance already exist", func() {
			JustBeforeEach(func() {
				response = makeRequest("DELETE", server.URL+"/instance/123456789", nil)
			})
			It("Delete the instance successfully", func() {
				Expect(response.StatusCode).Should(Equal(200))
			})
		})
	})
})
