package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	c, err := ReadConfig()
	if err != nil {
		t.Fail()
		t.Logf("Error found: %s", err.Error())
	}

	assert.Equal(t, "mqtt", c.Mqtt.Username, "MQTT username wrong.")
	assert.Equal(t, "heating", c.Mqtt.TopicPrefix, "MQTT prefix wrong.")
}
