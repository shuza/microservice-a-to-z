package client

import "github.com/shuza/microservice-a-to-z/vehicle-position-tracker/model"

/**
 * :=  created by:  Shuza
 * :=  create date:  24-Mar-2019
 * :=  (C) CopyRight Shuza
 * :=  www.shuza.ninja
 * :=  shuza.sa@gmail.com
 * :=  Fun  :  Coffee  :  Code
 **/

type IDbClient interface {
	Open() error
	UpdatePosition(position model.VehiclePosition)
	GetLastPositionFor(vehicleName string) (model.VehiclePosition, error)
	Disconnect()
}
