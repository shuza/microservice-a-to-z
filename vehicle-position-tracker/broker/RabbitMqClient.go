package broker

import (
	"fmt"
	"github.com/shuza/microservice-a-to-z/vehicle-position-tracker/model"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

/**
 * :=  created by:  Shuza
 * :=  create date:  24-Mar-2019
 * :=  (C) CopyRight Shuza
 * :=  www.shuza.ninja
 * :=  shuza.sa@gmail.com
 * :=  Fun  :  Coffee  :  Code
 **/

type RabbitMqClient struct {
	HostUrl   string
	Port      string
	Username  string
	Password  string
	Channel   *amqp.Channel
	QueueName string
}

func (c *RabbitMqClient) Init() error {
	c.HostUrl = "localhost"
	c.Port = "5672"
	c.Username = "admin"
	c.Password = "admin"
	c.QueueName = "position"

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

func (c *RabbitMqClient) GetQueue() (amqp.Queue, error) {
	return c.Channel.QueueDeclare(
		c.QueueName, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
}

func (c *RabbitMqClient) SendVehiclePosition(position model.VehiclePosition) error {
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
func (c *RabbitMqClient) GetConsumer() (<-chan amqp.Delivery, error) {
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

func failOnError(tag string, err error) {
	if err != nil {
		log.Fatalf("%s: %s", tag, err)
	}
}
