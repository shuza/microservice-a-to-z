package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

/**
 * :=  created by:  Shuza
 * :=  create date:  21-Mar-2019
 * :=  (C) CopyRight Shuza
 * :=  www.shuza.ninja
 * :=  shuza.sa@gmail.com
 * :=  Fun  :  Coffee  :  Code
 **/

func main() {
	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	failOnError(err, "Failed to connect RabbitMQ")
	defer conn.Close()

	log.Info("Connected successfully")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Errorf("%s\t==/\t%s", msg, err)
	}
}
