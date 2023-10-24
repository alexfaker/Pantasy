package config

import (
	"io/ioutil"
	"log"
	"os"
)

var Instance *Config //配置文件选项

// HTTPRequestTimeout ...
const HTTPRequestTimeout = 120
const AvatarNumberMax = 20

type EnvParams struct {
	UseAES                    bool   `yaml:"use_aes"`
	LogLevel                  string `yaml:"log_level"`
	CronTaskSwitch            bool   `yaml:"cron_task_switch"`
	DingTalkWebHookURL        string `yaml:"ding_talk_web_hook_url"`
	GORMLogLevel              int    `yaml:"gorm_log_level"`
	ASRAppID                  string `yaml:"asr_app_id"`
	ASRSecretID               string `yaml:"asr_secret_id"`
	ASRSecretKey              string `yaml:"asr_secret_key"`
	RedisFirstUseDB           int    `yaml:"redis_first_use_db"`
	SharePhotoQRCodeURL       string `yaml:"share_photo_qr_code_url"`
	AppPersonalPageURL        string `yaml:"app_personal_page_url"`
	NotificationTriggerSwitch bool   `yaml:"notification_trigger_switch"`
	IosAuditMode              bool   `yaml:"ios_audit_mode"`
	ForwardToAppServerAddress string `yaml:"forward_to_app_server_address"`
}

type Config struct {
	EnvParams EnvParams                            `yaml:"env_params"`
	Mysql     map[DatabaseBusiness]*databaseConfig `yaml:"mysql"` //数据库配置
	Redis     map[RedisBusiness]*redisConfig       `yaml:"redis"` //缓存配置
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
