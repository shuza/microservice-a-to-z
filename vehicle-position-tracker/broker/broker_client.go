package broker

import (
	"github.com/shuza/microservice-a-to-z/vehicle-position-tracker/model"
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

type IBrokerClient interface {
	Init() error
	GetQueue() (amqp.Queue, error)
	SendVehiclePosition(position model.VehiclePosition) error
	GetConsumer()
}
