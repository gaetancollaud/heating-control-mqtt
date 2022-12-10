package modules

import (
	"github.com/gaetancollaud/heating-control-mqtt/pkg/config"
	"github.com/gaetancollaud/heating-control-mqtt/pkg/data"
	"github.com/gaetancollaud/heating-control-mqtt/pkg/mqtt"
)

type HeatingModule struct {
	mqttClient    mqtt.Client
	heatingConfig []data.HeatingConfig
}

func (c *HeatingModule) Start() error {

	// TODO subscribe to all the topics
	return nil
}

func (c *HeatingModule) Stop() error {
	// TODO unsubscribe from all
	return nil
}

func NewHeatingModule(mqttClient mqtt.Client, heatingConfig []data.HeatingConfig, config *config.Config) Module {
	return &HeatingModule{
		mqttClient:    mqttClient,
		heatingConfig: heatingConfig,
	}
}

func init() {
	Register("heating", NewHeatingModule)
}
