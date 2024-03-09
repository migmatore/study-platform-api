package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Server   ServerConfig
	Logger   LoggerConfig
	Postgres PostgresConfig
}

type ServerConfig struct {
	RESTPort              string `mapstructure:"rest_port"`
	WSPort                string `mapstructure:"ws_port"`
	JwtSecretKey          string `mapstructure:"jwt_secret_key"`
	JwtExpTimeMin         int    `mapstructure:"jwt_exp_time_min"`
	WSJwtSecretKey        string `mapstructure:"ws_jwt_secret_key"`
	WSJwtExpTimeHour      int    `mapstructure:"ws_jwt_exp_time_hour"`
	JwtRefreshSecretKey   string `mapstructure:"jwt_refresh_secret_key"`
	JwtRefreshExpTimeHour int    `mapstructure:"jwt_refresh_exp_time_Hour"`
	Mode                  string `mapstructure:"mode"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name"`
}

func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}
