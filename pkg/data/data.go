package data

import "time"

type HeatingConfig struct {
	Name                   string
	TemperatureStatusTopic string
	PwmCommandTopic        string
	PwmStatusTopic         string
	OutputCommandTopic     string
	PwmRatio               float32
	PwmDutyCycle           time.Duration
}

func NewHeatingConfig() *HeatingConfig {
	return &HeatingConfig{
		Name:                   "",
		TemperatureStatusTopic: "",
		PwmCommandTopic:        "",
		PwmStatusTopic:         "",
		OutputCommandTopic:     "",
		PwmRatio:               0.5,
		PwmDutyCycle:           1 * time.Hour,
	}
}

func (o *HeatingConfig) SetName(value string) *HeatingConfig {
	o.Name = value
	return o
}

func (o *HeatingConfig) SetTemperatureStatusTopic(value string) *HeatingConfig {
	o.TemperatureStatusTopic = value
	return o
}

func (o *HeatingConfig) SetPwmCommandTopic(value string) *HeatingConfig {
	o.PwmCommandTopic = value
	return o
}

func (o *HeatingConfig) SetPwmStatusTopic(value string) *HeatingConfig {
	o.PwmStatusTopic = value
	return o
}

func (o *HeatingConfig) SetOutputCommandTopic(value string) *HeatingConfig {
	o.OutputCommandTopic = value
	return o
}

func (o *HeatingConfig) SetPwmRatio(value float32) *HeatingConfig {
	o.PwmRatio = value
	return o
}

func (o *HeatingConfig) SetPwmDutyCycle(value time.Duration) *HeatingConfig {
	o.PwmDutyCycle = value
	return o
}
