package brokerintegration_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/bcshuai/cf-redis-broker/integration"
)

var _ = Describe("Broker Security", func() {

	Describe("Basic HTTP Authentication", func() {

		Context("With expected username and password", func() {
			It("returns HTTP code 200", func() {
				code, _ := brokerClient.MakeCatalogRequest()
				Ω(code).To(Equal(200))
			})
		})

		Context("With unexpected username and password", func() {
			It("returns HTTP code 401", func() {
				client := &http.Client{}
				resp, err := client.Get("http://localhost:3000")
				defer resp.Body.Close()
				Ω(err).ToNot(HaveOccurred())
				req, err := http.NewRequest("GET", "http://localhost:3000/v2/catalog", nil)
				Ω(err).ToNot(HaveOccurred())
				req.SetBasicAuth("admin", "badpassword")
				resp, err = client.Do(req)
				Ω(err).ToNot(HaveOccurred())

				Ω(resp.StatusCode).To(Equal(401))
			})
		})
	})

	Describe("PORT Connectivity", func() {
		It("is available on the localhost only", func() {
			client := &http.Client{}
			resp, err := client.Get("http://localhost:3000")
			resp.Body.Close()
			Ω(err).ToNot(HaveOccurred())

			publicIPAddresses, err := integration.HostIP4Addresses()
			Ω(err).ToNot(HaveOccurred())

			for _, ipAddress := range publicIPAddresses {
				resp, err = client.Get(fmt.Sprintf("http://%s:3000", ipAddress))
				Ω(err).To(HaveOccurred())
			}
		})
	})
})
