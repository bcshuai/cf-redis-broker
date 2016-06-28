package agentconfig

import (
	"errors"
	"fmt"
	"os"

	"github.com/bcshuai/cf-redis-broker/brokerconfig"
	"github.com/cloudfoundry-incubator/candiedyaml"
)

type SharedAgentConfig struct {
	Host              string                         `yaml:"backend_host"`
	Port              string                         `yaml:"backend_port"`
	AuthConfiguration brokerconfig.AuthConfiguration `yaml:"auth"`

	ProcessCheckIntervalSeconds int `yaml:"process_check_interval"`
	StartRedisTimeoutSeconds    int `yaml:"start_redis_timeout"`
	ServiceInstanceLimit        int `yaml:"service_instance_limit"`

	RedisServerExecutablePath string `yaml:"redis_server_executable_path"`

	RedisConfiguration RedisConfiguration `yaml:"redis"`
}

type RedisConfiguration struct {
	DefaultConfigPath     string `yaml:"redis_conf_path"`
	InstanceDataDirectory string `yaml:"data_directory"`
	InstanceLogDirectory  string `yaml:"log_directory"`
}

func ParseSharedAgentConfig(path string) (SharedAgentConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return SharedAgentConfig{}, err
	}

	var config SharedAgentConfig
	if err := candiedyaml.NewDecoder(file).Decode(&config); err != nil {
		return SharedAgentConfig{}, err
	}

	return config, ValidateConfig(config.RedisConfiguration)
}

func ValidateConfig(config RedisConfiguration) error {
	err := checkPathExists(config.DefaultConfigPath, "RedisConfig.DefaultRedisConfPath")
	if err != nil {
		return err
	}

	err = checkPathExists(config.InstanceDataDirectory, "RedisConfig.InstanceDataDirectory")
	if err != nil {
		return err
	}

	err = checkPathExists(config.InstanceLogDirectory, "RedisConfig.InstanceLogDirectory")
	if err != nil {
		return err
	}

	return nil
}

func checkPathExists(path string, description string) error {
	_, err := os.Stat(path)
	if err != nil {
		errMessage := fmt.Sprintf(
			"File '%s' (%s) not found",
			path,
			description)
		return errors.New(errMessage)
	}
	return nil
}
