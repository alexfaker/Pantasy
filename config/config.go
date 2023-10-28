package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

var Instance *Config //配置文件选项

// HTTPRequestTimeout ...
const HTTPRequestTimeout = 120

type EnvironmentVariables struct {
	LogLevel     string `yaml:"loglevel"`
	GORMLogLevel int    `yaml:"gormLogLevel"`
}

type Config struct {
	Env   EnvironmentVariables                 `yaml:"environmentVariables"`
	Mysql map[DatabaseBusiness]*databaseConfig `yaml:"mysql"` //数据库配置
	Redis map[RedisBusiness]*redisConfig       `yaml:"redis"` //缓存配置
	Kafka map[RedisBusiness]*KafkaConfig       `yaml:"kafka"` //kafka配置
}

func Initialize() {
	//读取配置文件
	if err := LoadConfig(); err != nil {
		log.Fatal(err)
	}

	//初始化所有缓存
	for k, kVal := range Instance.Redis {
		localRedisManager.add(k, kVal.Address, kVal.Password)
	}
	localRedisManager.Ping()

	//初始化所有数据库
	for k, kVal := range Instance.Mysql {
		localDatabaseManager.add(k, kVal)
	}

}

func LoadConfig() error {
	vEnvironment := os.Getenv("env")
	defaultFile := "./conf/default.yaml"
	if vEnvironment == "local" {
		defaultFile = "./conf/local.yaml"
	}

	file, err := os.Open(defaultFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	c := &Config{}
	if err := yaml.Unmarshal(content, c); err != nil {
		return err
	}
	Instance = c
	return nil
}
