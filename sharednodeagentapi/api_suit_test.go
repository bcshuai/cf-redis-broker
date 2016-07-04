package sharednodeagentapi

import (
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("junit_api.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Shared Agent API Suite", []Reporter{junitReporter})
}
