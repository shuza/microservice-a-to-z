package main

import (
	"github.com/joho/godotenv"
	"github.com/shuza/microservice-a-to-z/vehicle-position-tracker/broker"
	"github.com/shuza/microservice-a-to-z/vehicle-position-tracker/client"
	log "github.com/sirupsen/logrus"
)

/**
 * :=  created by:  Shuza
 * :=  create date:  24-Mar-2019
 * :=  (C) CopyRight Shuza
 * :=  www.shuza.ninja
 * :=  shuza.sa@gmail.com
 * :=  Fun  :  Coffee  :  Code
 **/

func main() {
	godotenv.Load()
	initMongoDbClient()
	defer client.MongoDbClient.Disconnect()

	initMessageBroker()
	go broker.RabbitMqClient.ObserveAndStore(client.MongoDbClient)

}

func initMongoDbClient() {
	client.MongoDbClient = &client.MongoDbConnection{}
	err := client.MongoDbClient.Open()
	failOnError("Failed to Connect MongoDB", err)
}

func initMessageBroker() {
	broker.RabbitMqClient = &broker.RabbitMqConnection{}
	err := broker.RabbitMqClient.Init()
	failOnError("Client init failed", err)
}

func failOnError(tag string, err error) {
	if err != nil {
		log.Fatalf("%s: %s", tag, err)
	}
}
