package main

import (
	"encoding/json"
	"github.com/shuza/microservice-a-to-z/vehicle-position-tracker/broker"
	"github.com/shuza/microservice-a-to-z/vehicle-position-tracker/client"
	"github.com/shuza/microservice-a-to-z/vehicle-position-tracker/model"
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
	messageClient := broker.RabbitMqClient{}
	if err := messageClient.Init(); err != nil {
		failOnError("Client init failed", err)
	}

	dataChannel, err := messageClient.GetConsumer()
	failOnError("Failed to get position", err)

	mongoClient := client.MongoDbClient{}
	err = mongoClient.Open()
	failOnError("Failed to Connect MongoDB", err)

	for byteData := range dataChannel {
		data := model.VehiclePosition{}
		err := json.Unmarshal(byteData.Body, &data)
		if err != nil {
			log.Warnf("Json Parse Error : \t%s", err)
		} else {
			mongoClient.UpdatePosition(data)
		}
	}
}

func failOnError(tag string, err error) {
	if err != nil {
		log.Fatalf("%s: %s", tag, err)
	}
}
