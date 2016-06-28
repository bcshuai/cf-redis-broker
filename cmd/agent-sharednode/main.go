package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/bcshuai/brokerapi/auth"
	"github.com/bcshuai/cf-redis-broker/agentconfig"
	"github.com/bcshuai/cf-redis-broker/availability"
	"github.com/bcshuai/cf-redis-broker/process"
	"github.com/bcshuai/cf-redis-broker/redis"
	sharednode "github.com/bcshuai/cf-redis-broker/sharednodeagent"
	api "github.com/bcshuai/cf-redis-broker/sharednodeagentapi"
	"github.com/bcshuai/cf-redis-broker/system"

	"github.com/pivotal-golang/lager"
)

func main() {
	configPath := flag.String("agentConfig", "", "Shared node agent config yaml")
	flag.Parse()
	logger := lager.NewLogger("shared-node-redis-agent")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))

	//config, err := brokerconfig.ParseConfig(*configPath)
	config, err := agentconfig.ParseSharedAgentConfig(*configPath)
	if err != nil {
		logger.Fatal("Loading config file", err, lager.Data{
			"broker-config-path": configPath,
		})
	}

	commandRunner := system.OSCommandRunner{
		Logger: logger,
	}
	localRepo := &redis.LocalRepository{
		RedisConf: config,
	}
	processController := &redis.OSProcessController{
		CommandRunner:            commandRunner,
		InstanceInformer:         localRepo,
		Logger:                   logger,
		ProcessChecker:           &process.ProcessChecker{},
		ProcessKiller:            &process.ProcessKiller{},
		WaitUntilConnectableFunc: availability.Check,
	}
	localCreator := &redis.LocalInstanceCreator{
		FindFreePort:            system.FindFreePort,
		RedisConfiguration:      config,
		ProcessController:       processController,
		LocalInstanceRepository: localRepo,
	}
	agent := &sharednode.SharedNodeAgent{
		InstanceCreator: localCreator,
		InstanceRepo:    localCreator,
		Config:          config,
		Logger:          logger,
	}

	apiProvider := api.New(agent, logger)

	handler := auth.NewWrapper(
		config.AuthConfiguration.Username,
		config.AuthConfiguration.Password,
	).Wrap(
		apiProvider,
	)

	http.Handle("/", handler)
	logger.Fatal("http-listen", http.ListenAndServe(config.Host+":"+config.Port, nil))
}
