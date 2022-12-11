package modules

import (
	"github.com/gaetancollaud/heating-control-mqtt/pkg/config"
	"github.com/gaetancollaud/heating-control-mqtt/pkg/data"
	"github.com/gaetancollaud/heating-control-mqtt/pkg/mqtt"
	"github.com/gaetancollaud/heating-control-mqtt/pkg/pwm"
	"github.com/rs/zerolog/log"
)

type HeatingModule struct {
	mqttClient    mqtt.Client
	heatingConfig []data.HeatingConfig
}

func (c *HeatingModule) On(id string, data data.HeatingConfig) {
	log.Info().Str("id", id).Msg("ON")
	c.mqttClient.Publish(data.OutputCommandTopic, "{\"id\":1,\"src\":\"heating-control\",\"method\":\"Switch.Set\",\"params\": {\"id\": "+data.SwitchId+",\"on\": true}")
}

func (c *HeatingModule) Off(id string, data data.HeatingConfig) {
	log.Info().Str("id", id).Msg("OFF")
	c.mqttClient.Publish(data.OutputCommandTopic, "{\"id\":1,\"src\":\"heating-control\",\"method\":\"Switch.Set\",\"params\": {\"id\": "+data.SwitchId+",\"on\": off}")
}

func (c *HeatingModule) Start() error {

	// TODO subscribe to all the topics

	for _, heatingConfig := range c.heatingConfig {
		newPwm := pwm.NewPwm(heatingConfig.Name, heatingConfig.PwmDutyCycle, c, heatingConfig)

		// TODO listen from topic and restore
		newPwm.SetValuePercent(heatingConfig.PwmPercent)

		newPwm.Start()
	}

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
