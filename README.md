# prom2mqtt

Prom2MQTT is a tool that acts as a bridge from [Prometheus](https://prometheus.io/) to [MQTT](https://mqtt.org/).

## Installing Prom2MQTT

### Using `go install`

```console
go install github.com/suprememoocow/prom2mqtt@latest
```

## Running Prom2MQTT

Prom2MQTT has three required arguments:

1. `--config <file.yml>`: reference the configuration file to use. See the [Configuring Prom2MQTT](#configuring-prom2mqtt) section for more details.
1. `--mqtt.broker tcp://192.168.1.5:1883`: the MQTT broker to connect to.
1. `--prometheus.url http://192.168.1.5:9090`: the Prometheus service to use.

Putting these together, here's an example of running the command:

```console
$ prom2mqtt --config config.yaml --mqtt.broker tcp://192.168.138.5:1883 --prometheus.url http://192.168.138.5:9090
```

## Configuring Prom2MQTT

This example

```yaml
groups:
  # Multiple groups can be defined
  - name: weather_station

    # The interval between prometheus queries
    interval: 1m

    # Each group consists of 1 of more queries, each of which
    # will be published to a topic
    queries:
      # Publish the rain amount of rain over the past hour
      # to the weather/rain_1h topic
      - topic: weather/rain_1h
        expr: sum(increase(weather_station_rain_mm[1h]))
        decimal_places: 1

      # Publish the rain amount of rain over the 24 hours
      # to the weather/rain_24h topic
      - topic: weather/rain_24h
        expr: sum(increase(weather_station_rain_mm[24h]))
        decimal_places: 1

      # Publish a boolean 0/1 value which signals whether it's
      # rained in the past 30m to the weather/rained_last_30m topic
      - topic: weather/rained_last_30m
        expr: sum(increase(weather_station_rain_mm[30m])) > bool 0
        decimal_places: 0
```
