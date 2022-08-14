package group

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"

	"github.com/suprememoocow/prom2mqtt/internal/config"
)

var errUnexpectedResult = errors.New("unexpected result type from Prometheus result")

// Runner runs a group of queries on a fixed interval.
type Runner struct {
	v1api      v1.API
	mqttClient mqtt.Client
	group      *config.Group
}

// StartRunner will setup a new runner. This is intended to be
// called in a new go-routine.
func StartRunner(mqttBroker string, v1api v1.API, group *config.Group) {
	runner := &Runner{}
	runner.group = group
	runner.v1api = v1api

	opts := mqtt.NewClientOptions().AddBroker(mqttBroker).SetClientID("prom2mqtt")
	opts.SetKeepAlive(2 * time.Second)
	// opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)
	opts.AutoReconnect = true
	opts.ConnectRetry = true

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("mqtt connection error: %v", token.Error())
	}

	runner.mqttClient = c
	runner.run()
}

func (c *Runner) run() {
	for {
		startTime := time.Now()

		for _, q := range c.group.Queries {
			err := c.runQuery(q)
			if err != nil {
				fmt.Printf("query failed: %v", err)
			}
		}

		runTime := time.Since(startTime)

		remaining := c.group.Interval - runTime
		if remaining > 0 {
			time.Sleep(remaining)
		}
	}
}

// runQuery runs a single query as defined in the config and
// publishes it to MQTT.
func (c *Runner) runQuery(query *config.Query) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Query Prometheus using the supplied expression
	result, _, err := c.v1api.Query(ctx, query.Expr, time.Now(), v1.WithTimeout(5*time.Second))
	if err != nil {
		return fmt.Errorf("prometheus query failed: %w", err)
	}

	// Convert the Prometheus into an MQTT payload
	payload, err := c.resultToPayload(query, result)
	if err != nil {
		return fmt.Errorf("unable to create payload: %w", err)
	}

	c.mqttClient.Publish(query.Topic, 0, true, payload)

	return nil
}

func (c *Runner) resultToPayload(query *config.Query, result model.Value) (string, error) {
	vec, ok := result.(model.Vector)
	if !ok {
		return "", fmt.Errorf("result is not a vector: %w", errUnexpectedResult)
	}

	if len(vec) < 1 {
		return "", fmt.Errorf("empty vector: %w", errUnexpectedResult)
	}

	if len(vec) > 1 {
		log.Printf("more than one value, selecting the first")
	}

	return fmt.Sprintf("%f", vec[0].Value), nil
}
