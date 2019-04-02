package model

import (
	"encoding/json"
	"time"
)

/**
 * :=  created by:  Shuza
 * :=  create date:  24-Mar-2019
 * :=  (C) CopyRight Shuza
 * :=  www.shuza.ninja
 * :=  shuza.sa@gmail.com
 * :=  Fun  :  Coffee  :  Code
 **/

type VehiclePosition struct {
	Name      string    `json:"name"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`
}

func (model *VehiclePosition) ToJson() []byte {
	bytes, _ := json.Marshal(model)
	return bytes
}
