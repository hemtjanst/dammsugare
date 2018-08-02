package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hemtjanst/bibliotek/client"
	"github.com/hemtjanst/bibliotek/device"
	"github.com/hemtjanst/bibliotek/feature"
	"github.com/hemtjanst/bibliotek/transport/mqtt"
)

var (
	timeout      = flag.Int("timeout", 120, "Minutes after which we think the robot is done")
	startTopic   = flag.String("start.topic", "remote/robovac/KEY_AUTO", "LIRCd topic to start the vacuum")
	startPress   = flag.String("start.press", "200", "Milliseconds to hold down the start button")
	stopTopic    = flag.String("stop.topic", "remote/robovac/KEY_HOME", "LIRCd topic to stop the vacuum")
	stopPress    = flag.String("stop.press", "5000", "Milliseconds to hold down the stop button")
	manufacturer = flag.String("robot.manufacturer", "Eufy", "Vacuum manufacturer")
	name         = flag.String("robot.name", "RoboVac", "Vacuum name")
	model        = flag.String("robot.model", "11", "Vacuum model")
	serial       = flag.String("robot.serial-number", "undefined", "Vacuum serial number")
)

type handler struct {
	devices []*device.Device
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Parameters:\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}
	mcfg := mqtt.MustFlags(flag.String, flag.Bool)
	flag.Parse()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	m, err := mqtt.New(ctx, mcfg())
	if err != nil {
		panic(err)
	}

	robot, _ := client.NewDevice(&device.Info{
		Topic:        "robot/vacuum",
		Name:         *name,
		Model:        *model,
		Type:         "switch",
		Manufacturer: *manufacturer,
		SerialNumber: *serial,
		Features: map[string]*feature.Info{
			"on": &feature.Info{}},
	}, m)

	ft := robot.Feature("on")
	on, _ := ft.OnSet()
	var t *time.Timer

loop:
	for {
		select {
		case sig := <-quit:
			log.Printf("Received signal: %s, proceeding to shutdown", sig)
			break loop
		// Publish after every interval has elapsed
		case msg, open := <-on:
			if !open {
				break loop
			}
			switch msg {
			case "1":
				if t != nil {
					t.Stop()
				}
				m.Publish(*startTopic, []byte(*startPress), false)
				ft.Update("1")
				log.Print("Turned on robot")
				t = time.AfterFunc(time.Duration(*timeout)*time.Minute,
					func() {
						ft.Update("0")
						log.Print("Timeout expired, setting switch to off")
					})
			default:
				if t != nil {
					t.Stop()
				}
				m.Publish(*stopTopic, []byte(*stopPress), false)
				ft.Update("0")
				log.Print("Turned off robot")
			}
		}
	}
	cancel()
	os.Exit(0)
}
