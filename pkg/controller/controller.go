package controller

import (
	"fmt"
	"github.com/gaetancollaud/heating-control-mqtt/pkg/data"
	"time"

	"github.com/gaetancollaud/heating-control-mqtt/pkg/config"
	"github.com/gaetancollaud/heating-control-mqtt/pkg/controller/modules"
	"github.com/gaetancollaud/heating-control-mqtt/pkg/mqtt"
	"github.com/rs/zerolog/log"
)

type Controller struct {
	mqttClient    mqtt.Client
	heatingConfig []data.HeatingConfig

	modules map[string]modules.Module
}

func NewController(config *config.Config) *Controller {
	// Create Heating config
	var heatingConfig []data.HeatingConfig
	heatingConfig = append(heatingConfig, *data.NewHeatingConfig().
		SetName("gaetan").
		SetOutputCommandTopic("heating/gaetan/output").
		SetPwmCommandTopic("heating/gaetan/command").
		SetPwmDutyCycle(10 * time.Second))

	mqttOptions := mqtt.NewClientOptions().
		SetMqttUrl(config.Mqtt.MqttUrl).
		SetUsername(config.Mqtt.Username).
		SetPassword(config.Mqtt.Password).
		SetTopicPrefix(config.Mqtt.TopicPrefix).
		SetRetain(config.Mqtt.Retain)
	mqttClient := mqtt.NewClient(mqttOptions)

	controller := Controller{
		heatingConfig: heatingConfig,
		mqttClient:    mqttClient,
		modules:       map[string]modules.Module{},
	}

	for name, builder := range modules.Modules {
		module := builder(mqttClient, heatingConfig, config)
		controller.modules[name] = module
	}

	return &controller
}

func (c *Controller) Start() error {
	log.Info().Msg("Starting controller.")
	if err := c.mqttClient.Connect(); err != nil {
		return fmt.Errorf("error connecting to MQTT client: %w", err)
	}

	for name, module := range c.modules {
		log.Info().Str("module", name).Msg("Starting module.")
		if err := module.Start(); err != nil {
			return fmt.Errorf("error starting module '%s': %w", name, err)
		}
	}

	return nil
}

func (c *Controller) Stop() error {
	log.Info().Msg("Stopping controller.")

	for name, module := range c.modules {
		log.Info().Str("module", name).Msg("Stopping module.")
		if err := module.Stop(); err != nil {
			return fmt.Errorf("error stopping module '%s': %w", name, err)
		}
	}

	if err := c.mqttClient.Disconnect(); err != nil {
		return fmt.Errorf("error disconnecting to MQTT client: %w", err)
	}

	return nil
}
