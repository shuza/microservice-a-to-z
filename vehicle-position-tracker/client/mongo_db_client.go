package client

import (
	"context"
	"fmt"
	"github.com/shuza/microservice-a-to-z/vehicle-position-tracker/model"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
type MongoDbClient struct {
	hostUrl string
	port    string
	client  *mongo.Client
}

func (c *MongoDbClient) Open() error {
	c.hostUrl = "localhost"
	c.port = "27017"
	connectionStr := fmt.Sprintf("mongodb://%s:%s", c.hostUrl, c.port)
	clientOptions := options.Client().ApplyURI(connectionStr)
	conn, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}
	c.client = conn

	//	Check the connection
	err = c.client.Ping(context.TODO(), nil)
	return err
}

func (c *MongoDbClient) UpdatePosition(position model.VehiclePosition) {
	collection := c.client.Database("test").Collection("vehicle")
	insertResult, err := collection.InsertOne(context.TODO(), position)

	if err != nil {
		log.Warnf("%s\t:\t%s", "Can't update "+position.Name, err)
	} else {
		log.Infof("%s\tupdated ID ==\t%s"+position.Name, insertResult)
	}
}

func (c *MongoDbClient) GetLastPositionFor(vehicleName string) (model.VehiclePosition, error) {
	cursor, err := c.client.Database("test").Collection("vehicle").Find(context.Background(), bson.D{})
	if err != nil {
		log.Warnf("Get last position for %s Error \t:\t%s", vehicleName, err)
		return model.VehiclePosition{}, err
	}
	defer cursor.Close(context.Background())
	position := model.VehiclePosition{}
	for cursor.Next(context.Background()) {
		err = cursor.Decode(&position)
	}

	return position, nil

}

func (c *MongoDbClient) Disconnect() {
	c.client.Disconnect(context.TODO())
}
