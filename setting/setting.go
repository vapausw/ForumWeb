package setting

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

var Conf = new(AppConfig)

type AppConfig struct {
	AppName        string `mapstructure:"name"`
	Port           int    `mapstructure:"port"`
	Version        string `mapstructure:"version"`
	Mode           string `mapstructure:"mode"`
	StartTime      string `mapstructure:"start_time"`
	MachineID      int64  `mapstructure:"machine_id"`
	*LogConfig     `mapstructure:"log"`
	*MySQLConfig   `mapstructure:"mysql"`
	*RedisConfig   `mapstructure:"redis"`
	*MyEmailConfig `mapstructure:"email"`
	*MyKafkaConfig `mapstructure:"kafka"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
}

type MySQLConfig struct {
	Username     string `mapstructure:"user_name"`
	Password     string `mapstructure:"password"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	DbName       string `mapstructure:"db_name"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	Db       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type MyEmailConfig struct {
	Email    string `mapstructure:"email"`
	Password string `mapstructure:"password"`
	SmtpHost string `mapstructure:"smtp_host"`
	SmtpPort string `mapstructure:"smtp_port"`
}

type MyKafkaConfig struct {
	Brokers    []string `mapstructure:"brokers"`
	EmailTopic string   `mapstructure:"email_topic"`
	GroupID    string   `mapstructure:"group_id"`
}

func Init(filepath string) (err error) {
	viper.SetConfigFile(filepath)
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig failed, err:%v\n", err)
		return
	}
	if err = viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		return
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config.yaml has been changed")
		if err = viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
			return
		}
	})
	return
}
