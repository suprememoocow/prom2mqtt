package main

import (
	"flag"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"

	"github.com/suprememoocow/prom2mqtt/internal/config"
	"github.com/suprememoocow/prom2mqtt/internal/group"
)

func main() {
	configFile := flag.String("config", "", "config file name")
	mqttBroker := flag.String("mqtt.broker", "", "mqtt broker - eg tcp://localhost:1883")
	prometheusEndpoint := flag.String("prometheus.url", "", "prometheus url - eg http://localhost:9090")
	flag.Parse()

	// MQTT logs to stderr
	mqtt.ERROR = log.New(os.Stderr, "", 0)

	if *configFile == "" {
		log.Fatalf("config file required")
	}

	conf, err := config.Load(*configFile)
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}

	client, err := api.NewClient(api.Config{
		Address: *prometheusEndpoint,
	})
	if err != nil {
		log.Fatalf("error creating client: %v\n", err)
	}

	v1api := v1.NewAPI(client)

	for _, v := range conf.Groups {
		go group.StartRunner(*mqttBroker, v1api, v)
	}

	select {}
}
