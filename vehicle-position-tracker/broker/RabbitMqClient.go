package broker

import (
	"encoding/json"
	"fmt"
	"github.com/shuza/microservice-a-to-z/vehicle-position-tracker/client"
	"github.com/shuza/microservice-a-to-z/vehicle-position-tracker/model"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
)

/**
 * :=  created by:  Shuza
 * :=  create date:  24-Mar-2019
 * :=  (C) CopyRight Shuza
 * :=  www.shuza.ninja
 * :=  shuza.sa@gmail.com
 * :=  Fun  :  Coffee  :  Code
 **/

type RabbitMqConnection struct {
	HostUrl   string
	Port      string
	Username  string
	Password  string
	Channel   *amqp.Channel
	QueueName string
}

var RabbitMqClient IBrokerClient

func (c *RabbitMqConnection) Init() error {
	c.HostUrl = os.Getenv("RABBIT_MQ_HOST")
	c.Port = os.Getenv("RABBIT_MQ_PORT")
	c.Username = os.Getenv("RABBIT_MQ_USERNAME")
	c.Password = os.Getenv("RABBIT_MQ_PASSWORD")
	c.QueueName = os.Getenv("RABBIT_MQ_QUEUE")

	connectionStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", c.Username, c.Password, c.HostUrl, c.Port)
	conn, err := amqp.Dial(connectionStr)
	if err != nil {
		conn.Close()
		return err
	}

	if c.Channel, err = conn.Channel(); err != nil {
		c.Channel.Close()
		return err
	}

	return nil
}

func (c *RabbitMqConnection) GetQueue() (amqp.Queue, error) {
	return c.Channel.QueueDeclare(
		c.QueueName, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
}

func (c *RabbitMqConnection) SendVehiclePosition(position model.VehiclePosition) error {
	err := c.Channel.Publish(
		"",
		c.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(position.ToJson()),
		},
	)
	return err
}

//	return channel for receiving data
func (c *RabbitMqConnection) GetConsumer() (<-chan amqp.Delivery, error) {
	return c.Channel.Consume(
		c.QueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (c *RabbitMqConnection) ObserveAndStore(dbClient client.IDbClient) {
	dataChannel, err := c.GetConsumer()
	failOnError("Failed to consume message", err)

	for byteData := range dataChannel {
		data := model.VehiclePosition{}
		if err := json.Unmarshal(byteData.Body, &data); err == nil {
			client.MongoDbClient.UpdatePosition(data)
		} else {
			log.Warnf("Json message Parse Error : \t%s", err)
		}
	}
}

func failOnError(tag string, err error) {
	if err != nil {
		log.Fatalf("%s: %s", tag, err)
	}
}
