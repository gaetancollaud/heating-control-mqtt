package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type ConfigMqtt struct {
	MqttUrl     string
	Username    string
	Password    string
	TopicPrefix string
	Retain      bool
}

type Config struct {
	Mqtt     ConfigMqtt
	LogLevel string
}

const (
	undefined             string = "__undefined__"
	deprecated            string = "__deprecated__"
	configFile            string = "config.yaml"
	envKeyMqttUrl         string = "mqtt_url"
	envKeyMqttUsername    string = "mqtt_username"
	envKeyMqttPassword    string = "mqtt_password"
	envKeyMqttTopicPrefix string = "mqtt_topic_prefix"
	envKeyMqttRetain      string = "mqtt_retain"
	envKeyLogLevel        string = "log_level"
)

var defaultConfig = map[string]interface{}{
	envKeyMqttUrl:         undefined,
	envKeyMqttUsername:    undefined,
	envKeyMqttPassword:    undefined,
	envKeyMqttTopicPrefix: "heating",
	envKeyMqttRetain:      false,
	envKeyLogLevel:        "INFO",
}

// FromEnv returns a Config from env variables
func ReadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	// Set the current directory where the binary is being run.
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	for key, value := range defaultConfig {
		if value != undefined && value != deprecated {
			viper.SetDefault(key, value)
		}
	}

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("ReadInConfig error: %w", err)
	}

	// Check for deprecated and undefined fields.
	for fieldName, defaultValue := range defaultConfig {
		if defaultValue == deprecated && viper.IsSet(fieldName) {
			return nil, fmt.Errorf("deprecated field found in config: %s", fieldName)
		}
		//if defaultValue == undefined && !viper.IsSet(fieldName) {
		//	return nil, fmt.Errorf("required field not found in config: %s", fieldName)
		//}
	}

	config := &Config{
		Mqtt: ConfigMqtt{
			MqttUrl:     viper.GetString(envKeyMqttUrl),
			Username:    viper.GetString(envKeyMqttUsername),
			Password:    viper.GetString(envKeyMqttPassword),
			TopicPrefix: viper.GetString(envKeyMqttTopicPrefix),
			Retain:      viper.GetBool(envKeyMqttRetain),
		},
		LogLevel: viper.GetString(envKeyLogLevel),
	}

	return config, nil
}

func (c *Config) String() string {
	return fmt.Sprintf("%+v\n", c.Mqtt)
}
