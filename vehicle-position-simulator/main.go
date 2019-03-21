package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/tealeg/xlsx"
	"microservice-a-to-z/vehicle-position-simulator/model"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
	//dummyDataDir := "/home/sanjay/go/src/microservice-a-to-z/vehicle-position-simulator/data"
	currentDir, _ := os.Getwd()
	dummyDataDir := currentDir + "/data/"

	channel, err := newChannel()
	failOnError(err, "Failed to connect Message Broker")

	queue, err := getQueue(channel)
	failOnError(err, "Failed to connect with Message Queue")

	var coroutineWaiter sync.WaitGroup
	err = filepath.Walk(dummyDataDir, func(path string, info os.FileInfo, err error) error {
		if xlFile, err := xlsx.OpenFile(path); err == nil {
			parts := strings.Split(info.Name(), ".")
			go simulateVehiclePosition(&coroutineWaiter, channel, queue, parts[0], xlFile)
			coroutineWaiter.Add(1)
		} else {
			warnOnError(err, "Can't read "+info.Name()+" file")
		}
		return nil
	})

	failOnError(err, "Failed to read Dir")
	coroutineWaiter.Wait()
	log.Info("\n\t====\tFinished\t====")
}

func simulateVehiclePosition(coroutineWaiter *sync.WaitGroup, channel *amqp.Channel, queue amqp.Queue, vehicleName string, xlFile *xlsx.File) {
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			if len(row.Cells) >= 2 {
				lat, errLat := row.Cells[0].Float()
				lng, errLng := row.Cells[1].Float()
				if errLat != nil || errLng != nil {
					continue
				}

				position := model.VehiclePosition{vehicleName, lat, lng}
				sendPosition(position, channel, queue)
			}
		}
	}

	coroutineWaiter.Done()
}

func sendPosition(position model.VehiclePosition, channel *amqp.Channel, queue amqp.Queue) {
	err := channel.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(position.ToJson()),
		},
	)
	warnOnError(err, "Failed to push position for "+position.Name)
	log.Infof("Update \t==\t%s", position.Name)
}

func newChannel() (*amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	if err != nil {
		return &amqp.Channel{}, err
	}

	return conn.Channel()
}

func getQueue(channel *amqp.Channel) (amqp.Queue, error) {
	return channel.QueueDeclare(
		"position", // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
}

/*func main() {
	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	failOnError(err, "Failed to connect RabbitMQ")
	defer conn.Close()

	log.Info("Connected successfully")
}*/

func failOnError(err error, msg string) {
	if err != nil {
		log.Errorf("%s\t==/\t%s", msg, err)
	}
}

func warnOnError(err error, msg string) {
	if err != nil {
		log.Warnf("%s\t==/\t%s", msg, err)
	}
}
