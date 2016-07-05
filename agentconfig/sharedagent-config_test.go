package agentconfig_test

import (
	"path"
	"path/filepath"

	"github.com/bcshuai/cf-redis-broker/agentconfig"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Check the shared node config", func() {
	Describe("Load the config file", func() {
		Context("the config file does not exist", func() {
			It("return an error", func() {
				_, err := agentconfig.ParseSharedAgentConfig("/some/invalid/path")
				Expect(err.Error()).To(Equal("open /some/invalid/path: no such file or directory"))
			})
		})
		Context("When the config file is loaded successfully", func() {
			var config agentconfig.SharedAgentConfig

			BeforeEach(func() {
				path, err := filepath.Abs(path.Join("assets", "shared-agent.yml"))
				Expect(err).ToNot(HaveOccurred())

				config, err = agentconfig.ParseSharedAgentConfig(path)
				Expect(err).ToNot(HaveOccurred())

			})

			It("Has the correct conf path", func() {
				Expect(config.RedisConfiguration.DefaultConfigPath).To(Equal("/usr/bin/true"))
			})

			It("Has the right service_instance_limit", func() {
				Expect(config.ServiceInstanceLimit).To(Equal(5))
			})

			It("Has the right host", func() {
				Expect(config.Host).To(Equal("localhost"))
			})

			It("Has the right port", func() {
				Expect(config.Port).To(Equal("8888"))
			})

			It("Has the right user name", func() {
				Expect(config.AuthConfiguration.Username).To(Equal("admin"))
			})

			It("Has the right user password", func() {
				Expect(config.AuthConfiguration.Password).To(Equal("admin"))
			})
		})
	})
})
