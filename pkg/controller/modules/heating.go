package modules

import (
	mqtt2 "github.com/eclipse/paho.mqtt.golang"
	"github.com/gaetancollaud/heating-control-mqtt/pkg/config"
	"github.com/gaetancollaud/heating-control-mqtt/pkg/data"
	"github.com/gaetancollaud/heating-control-mqtt/pkg/mqtt"
	"github.com/gaetancollaud/heating-control-mqtt/pkg/pwm"
	"github.com/gaetancollaud/heating-control-mqtt/pkg/utils"
	"github.com/rs/zerolog/log"
	"strconv"
)

type HeatingModule struct {
	mqttClient    mqtt.Client
	heatingConfig []data.HeatingConfig
}

func (c *HeatingModule) setValue(data data.HeatingConfig, on bool) {
	log.Info().Str("name", data.Name).Bool("on", on).Msg("heating setValue")
	onStr := strconv.FormatBool(on)
	err := c.mqttClient.Publish(data.OutputCommandTopic, "{\"id\":1,\"src\":\"heating-control\",\"method\":\"Switch.Set\",\"params\": {\"id\": "+data.SwitchId+",\"on\": "+onStr+"}")
	utils.CheckNoErrorAndPrint(err)
}

func (c *HeatingModule) parsePvmCommand(pwm *pwm.Pwm, data data.HeatingConfig, input string) {
	log.Info().Str("name", data.Name).Str("input", input).Msg("heating parsePwmValue")
	i, err := strconv.ParseInt(input, 10, 32)
	if err == nil {
		if i > 100 {
			i = 100
		} else if i < 0 {
			i = 0
		}
		percent := uint8(i)
		data.SetPwmPercent(percent)
		pwm.SetValuePercent(percent)

		// publish back the status
		err := c.mqttClient.PublishWithPrefix(data.PwmStatusTopic, strconv.Itoa(int(data.PwmPercent)))
		utils.CheckNoErrorAndPrint(err)
	} else {
		log.Error().Str("input", input).Msg("heating unable to parse percent")
	}
}

func (c *HeatingModule) Start() error {

	// TODO subscribe to all the topics

	for _, heatingConfig := range c.heatingConfig {
		pwm := pwm.NewPwm(heatingConfig.Name, heatingConfig.PwmDutyCycle, func(id string, status bool) {
			c.setValue(heatingConfig, status)
		})

		pwm.SetValuePercent(heatingConfig.PwmPercent)

		err := c.mqttClient.PublishWithPrefix(heatingConfig.PwmStatusTopic, strconv.Itoa(int(heatingConfig.PwmPercent)))
		utils.CheckNoErrorAndPrint(err)

		c.mqttClient.SubscribeWithPrefix(heatingConfig.PwmCommandTopic, func(client mqtt2.Client, message mqtt2.Message) {
			c.parsePvmCommand(pwm, heatingConfig, string(message.Payload()))
		})

		pwm.Start()
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
