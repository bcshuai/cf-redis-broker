package main

import (
	"net/http"
	"os"

	"github.com/bcshuai/brokerapi"
	"github.com/bcshuai/brokerapi/auth"
	"github.com/pivotal-golang/lager"

	"github.com/bcshuai/cf-redis-broker/availability"
	"github.com/bcshuai/cf-redis-broker/broker"
	"github.com/bcshuai/cf-redis-broker/brokerconfig"
	"github.com/bcshuai/cf-redis-broker/debug"
	"github.com/bcshuai/cf-redis-broker/process"
	"github.com/bcshuai/cf-redis-broker/redis"
	"github.com/bcshuai/cf-redis-broker/redis/sharednode"
	"github.com/bcshuai/cf-redis-broker/redisinstance"
	"github.com/bcshuai/cf-redis-broker/system"
)

func main() {
	brokerConfigPath := configPath()

	brokerLogger := lager.NewLogger("redis-broker")
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))

	brokerLogger.Info("Config File: " + brokerConfigPath)

	config, err := brokerconfig.ParseConfig(brokerConfigPath)
	if err != nil {
		brokerLogger.Fatal("Loading config file", err, lager.Data{
			"broker-config-path": brokerConfigPath,
		})
	}

	commandRunner := system.OSCommandRunner{
		Logger: brokerLogger,
	}

	localRepo := &redis.LocalRepository{
		RedisConf: config.RedisConfiguration,
	}

	processController := &redis.OSProcessController{
		CommandRunner:            commandRunner,
		InstanceInformer:         localRepo,
		Logger:                   brokerLogger,
		ProcessChecker:           &process.ProcessChecker{},
		ProcessKiller:            &process.ProcessKiller{},
		WaitUntilConnectableFunc: availability.Check,
	}

	localCreator := &redis.LocalInstanceCreator{
		FindFreePort:            system.FindFreePort,
		RedisConfiguration:      config.RedisConfiguration,
		ProcessController:       processController,
		LocalInstanceRepository: localRepo,
	}

	agentClient := &redis.RemoteAgentClient{
		HttpAuth: config.AuthConfiguration,
	}
	remoteRepo, err := redis.NewRemoteRepository(agentClient, config)
	if err != nil {
		brokerLogger.Fatal("Error initializing remote repository", err)
	}

	//added by Bin Cheng Shuai to support multi agent client
	sharedRemoteRepo, err := sharednode.NewShareRemoteRepository(config, brokerLogger)
	if err != nil {
		brokerLogger.Fatal("Error initializing shared remote repository", err)
	}
	serviceBroker := &broker.BluemixRedisServiceBroker{
		RedisServiceBroker: broker.RedisServiceBroker{
			InstanceCreators: map[string]broker.InstanceCreator{
				"shared":    sharedRemoteRepo, //localCreator,
				"dedicated": remoteRepo,
				"local":     localCreator,
			},
			InstanceBinders: map[string]broker.InstanceBinder{
				"shared":    sharedRemoteRepo, //localRepo,
				"dedicated": remoteRepo,
				"local":     localRepo,
			},
			Config: config,
		},
	}

	brokerCredentials := brokerapi.BrokerCredentials{
		Username: config.AuthConfiguration.Username,
		Password: config.AuthConfiguration.Password,
	}

	brokerAPI := brokerapi.New(serviceBroker, brokerLogger, brokerCredentials)

	authWrapper := auth.NewWrapper(brokerCredentials.Username, brokerCredentials.Password)
	debugHandler := authWrapper.WrapFunc(debug.NewHandler(remoteRepo))
	instanceHandler := authWrapper.WrapFunc(redisinstance.NewHandler(remoteRepo))

	http.HandleFunc("/instance", instanceHandler)
	http.HandleFunc("/debug", debugHandler)
	http.Handle("/", brokerAPI)

	brokerLogger.Fatal("http-listen", http.ListenAndServe(config.Host+":"+config.Port, nil))
}

func configPath() string {
	brokerConfigYamlPath := os.Getenv("BROKER_CONFIG_PATH")
	if brokerConfigYamlPath == "" {
		panic("BROKER_CONFIG_PATH not set")
	}
	return brokerConfigYamlPath
}
