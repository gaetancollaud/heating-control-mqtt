package mqtt

import (
	"fmt"
	"path"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const (
	Online  string = "online"
	Offline string = "offline"
)

// Topics.
const (
	State        string = "state"
	Command      string = "command"
	Event        string = "event"
	serverStatus string = "server/status"
)

type Client interface {
	// Connect to the MQTT server.
	Connect() error
	// Disconnect from the MQTT server.
	Disconnect() error

	// Publishes a message under the prefix topic
	PublishWithPrefix(topic string, message interface{}) error
	Publish(topic string, message interface{}) error
	// Subscribe to a topic and calls the given handler when a message is
	// received.
	SubscribeWithPrefix(topic string, messageHandler mqtt.MessageHandler) error

	// Return the full topic for a given subpath.
	GetFullTopic(topic string) string
	// Returns the topic used to publish the server status.
	ServerStatusTopic() string

	RawClient() mqtt.Client
}

type client struct {
	mqttClient mqtt.Client
	options    ClientOptions
}

func NewClient(options *ClientOptions) Client {
	mqttOptions := mqtt.NewClientOptions().
		AddBroker(options.MqttUrl).
		SetClientID("heating-mqtt-" + uuid.New().String()).
		SetOrderMatters(false).
		SetUsername(options.Username).
		SetPassword(options.Password)

	return &client{
		mqttClient: mqtt.NewClient(mqttOptions),
		options:    *options,
	}
}

func (c *client) Connect() error {
	t := c.mqttClient.Connect()
	<-t.Done()
	if t.Error() != nil {
		return fmt.Errorf("error connecting to MQTT broker: %w", t.Error())
	}

	if err := c.publishServerStatus(Online); err != nil {
		return err
	}
	return nil
}

func (c *client) Disconnect() error {
	log.Info().Msg("Publishing Offline status to MQTT server.")
	if err := c.publishServerStatus(Offline); err != nil {
		return err
	}
	c.mqttClient.Disconnect(uint(c.options.DisconnectTimeout.Milliseconds()))
	log.Info().Msg("Disconnected from MQTT server.")
	return nil
}

func (c *client) PublishWithPrefix(topic string, message interface{}) error {
	return c.Publish(path.Join(c.options.TopicPrefix, topic), message)
}

func (c *client) Publish(topic string, message interface{}) error {
	t := c.mqttClient.Publish(
		topic,
		c.options.QoS,
		c.options.Retain,
		message)
	<-t.Done()
	return t.Error()
}

func (c *client) SubscribeWithPrefix(topic string, messageHandler mqtt.MessageHandler) error {
	t := c.mqttClient.Subscribe(
		path.Join(c.options.TopicPrefix, topic),
		c.options.QoS,
		messageHandler)
	<-t.Done()
	return t.Error()
}

// Publish the current binary status into the MQTT topic.
func (c *client) publishServerStatus(message string) error {
	log.Info().Str("status", message).Str("topic", serverStatus).Msg("Updating server status topic")
	return c.PublishWithPrefix(serverStatus, message)
}

func (c *client) ServerStatusTopic() string {
	return path.Join(c.options.TopicPrefix, serverStatus)
}

func (c *client) GetFullTopic(topic string) string {
	return path.Join(c.options.TopicPrefix, topic)
}

func (c *client) RawClient() mqtt.Client {
	return c.mqttClient
}
