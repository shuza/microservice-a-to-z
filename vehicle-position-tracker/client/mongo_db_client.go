package client

import (
	"fmt"
	"github.com/shuza/microservice-a-to-z/vehicle-position-tracker/model"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/mgo.v2"
	"os"
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

//	Real implementation of mongoDB client
type MongoDbConnection struct {
	conn *mgo.Session
}

var MongoDbClient IDbClient

func (c *MongoDbConnection) Open() error {
	connectionStr := fmt.Sprintf("mongodb://%s:%s",
		os.Getenv("MONGO_HOST"), os.Getenv("MONGO_PORT"))
	if conn, err := mgo.Dial(connectionStr); err != nil {
		return err
	} else {
		c.conn = conn
	}

	return nil
}

func (c *MongoDbConnection) UpdatePosition(position model.VehiclePosition) {
	position.CreatedAt = time.Now()
	collection := c.conn.DB(os.Getenv("MONGO_DB_NAME")).C("vehicle")
	if err := collection.Insert(position); err != nil {
		log.Warnf("%s\t:\t%s", "Can't update "+position.Name, err)
	}
}

func (c *MongoDbConnection) GetLastPositionFor(vehicleName string) (model.VehiclePosition, error) {
	var positions []model.VehiclePosition
	err := c.conn.DB(os.Getenv("MONGO_DB_NAME")).C("vehicle").
		Find(bson.M{"name": vehicleName}).
		Sort("-created_at").
		All(&positions)

	if err != nil {
		return model.VehiclePosition{}, err
	}

	index := len(positions) - 1
	position := positions[index]
	return position, nil
}

func (c *MongoDbConnection) Disconnect() {
	c.conn.Close()
}
