package configmigratorintegration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/bcshuai/cf-redis-broker/integration/helpers"

	"testing"
)

func TestConfigmigratorintegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Configmigrator Integration Suite")
}

var _ = BeforeEach(func() {
	helpers.ResetTestDirs()
})
