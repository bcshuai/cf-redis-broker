package brokerintegration_test

import (
	"net/http"

	"code.google.com/p/go-uuid/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dedicated instance unbinding", func() {
	var instanceID string
	var bindingID string

	BeforeEach(func() {
		instanceID = uuid.NewRandom().String()
		bindingID = uuid.NewRandom().String()

		code, _ := brokerClient.ProvisionInstance(instanceID, "dedicated")
		Ω(code).Should(Equal(201))

		status, _ := bindInstance(instanceID, bindingID)
		Ω(status).Should(Equal(http.StatusCreated))
	})

	AfterEach(func() {
		deprovisionInstance(instanceID)
	})

	It("should respond correctly", func() {
		code, body := unbindInstance(instanceID, bindingID)
		Ω(code).Should(Equal(200))
		Ω(body).Should(MatchJSON("{}"))

		code, body = unbindInstance(instanceID, bindingID)
		Ω(code).To(Equal(410))
		Ω(body).Should(MatchJSON("{}"))

		code, body = unbindInstance("NON-EXISTANT", bindingID)
		Ω(code).To(Equal(404))
		Ω(body).Should(MatchJSON("{}"))
	})
})
