package config

import (
	"time"

	config "github.com/justcgh9/go-config"
)

type Config struct {
	Env            string     `yaml:"env" env-default:"local"`
	StoragePath    string     `yaml:"storage_path" env-required:"true"`
	GraphStorage   GraphStorage	`yaml:"graph_storage" env-required:"true"`
	GRPCSrv        GRPCServerConfig `yaml:"grpc_srv" env-required:"true"`
	Kafka 		   Kafka `yaml:"kafka" env-required:"true"`
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

type GraphStorage struct {
	URI 		string 	`yaml:"uri" env-required:"true"`
	Username	string	`yaml:"username" env-required:"true"`
	Password 	string 	`yaml:"password" env-required:"true"`
	Realm 		string 	`yaml:"realm" env-default:"true"`
}

type Kafka struct {
	Brokers		[]string `yaml:"brokers" env-required:"true"`
	GroupID 	string `yaml:"group_id" env-default:"friends-consumer"`
	MinBytes 	int `yaml:"min_bytes" env-default:"10"`
	MaxBytes 	int `yaml:"max_bytes" env-default:"10e6"`
	MaxWait		time.Duration `yaml:"max_wait" env-default:"3s"`
}

func MustLoad() *Config {
	return config.MustLoad[Config]()
}


