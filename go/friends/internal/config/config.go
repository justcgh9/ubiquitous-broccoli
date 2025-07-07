package config

import (
	"time"

	config "github.com/justcgh9/go-config"
)

type Config struct {
	Env            string     `yaml:"env" env-default:"local"`
	StoragePath    string     `yaml:"storage_path" env-required:"true"`
	GRPCSrv        GRPCServerConfig `yaml:"grpc_srv" env-required:"true"`
	UsersClient	   GRPCUserServiceClient `yaml:"users_grpc_client" env-required:"true"`
}

type GRPCServerConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type GRPCUserServiceClient struct {
	URI	   		string 		`yaml:"uri" env-required:"true"`
	Timeout		time.Duration 	`yaml:"timeout"`
	QueueSize 	int 	`yaml:"queue_size" env-default:"16"`
	NumWorkers 	int 	`yaml:"num_workers" env-default:"4"`
}

func MustLoad() *Config {
	return config.MustLoad[Config]()
}


