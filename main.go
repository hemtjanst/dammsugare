package main

import (
	"flag"
	"fmt"
	"github.com/hemtjanst/hemtjanst/device"
	"github.com/hemtjanst/hemtjanst/messaging"
	"github.com/hemtjanst/hemtjanst/messaging/flagmqtt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	timeout = flag.Int("timeout", 120, "Minutes after which we think the robot is done")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Parameters:\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}
	flag.Parse()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	id := flagmqtt.NewUniqueIdentifier()
	conf := flagmqtt.ClientConfig{
		ClientID:    "dammsugare",
		WillTopic:   "leave",
		WillPayload: id,
		WillRetain:  false,
		WillQoS:     0,
	}
	c, err := flagmqtt.NewPersistentMqtt(conf)
	if err != nil {
		log.Fatal("Could not configure the MQTT client: ", err)
	}

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("Failed to establish connection with broker: ", token.Error())
	}

	m := messaging.NewMQTTMessenger(c)

	robot := device.NewDevice("robot/vacuum", m)
	robot.Manufacturer = "Eufy"
	robot.Name = "RoboVac"
	robot.Model = "11"
	robot.Type = "switch"
	robot.LastWillID = id
	robot.Features = map[string]*device.Feature{
		"on": {},
	}

	robot.PublishMeta()
	m.Subscribe("discover", 1, func(msg messaging.Message) {
		m.Publish("announce", []byte(robot.Topic), 1, false)
	})

	log.Print("Announced RoboVac to Hemtjänst")

	robot.OnSet("on", func(msg messaging.Message) {
		if string(msg.Payload()) == "1" {
			m.Publish("remote/robovac/KEY_AUTO", []byte("200"), 1, false)
			robot.Update("on", "1")
			go func() {
				<-time.After(time.Duration(*timeout) * time.Minute)
				robot.Update("on", "0")
			}()
			log.Print("Turned on robot")
		} else {
			m.Publish("remote/robovac/KEY_HOME", []byte("5000"), 1, false)
			robot.Update("on", "0")
			log.Print("Turned off robot")
		}
	})

loop:
	for {
		select {
		case sig := <-quit:
			log.Printf("Received signal: %s, proceeding to shutdown", sig)
			break loop
		}
	}

	c.Disconnect(250)
	log.Print("Disconnected from broker. Bye!")
	os.Exit(0)
}
