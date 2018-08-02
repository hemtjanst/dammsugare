# Dammsugare

Dammsugare exposes a switch to HomeKit that allows turning a robot vacuum
on or off. By default it is configured for a Eufy RoboVac 11.

**Note**: This relies on [rodljus][lirc]. The RoboVac is fairly dumb, it
doesn't have WiFi or other cloud-y API things to control it with. So instead
it is controlled by sending the necessary commands to mqtt2lirc.

## Installation

The installation is pretty simple, `go install` it and run it. Potentially
adjust the timeout and set a different manufacturer, name, model and serial
number.

```
go install github.com/hemtjanst/dammsugare
dammsugare -timeout 100 -mqtt.address broker.mydomain.tld:1883
```

Please note that using go install will grab the latest version of all the dependencies which can cause problems. As such a `Gopkg.toml` and associated lock file is provided to be used with dep ensure. After that you can `go build -o dammsugare main.go`.

## Configuration

Because the robot cannot signal back to us when it's done the timeout is used
to flip the `on` state back to off. Try to guestimate how long the cleaning
cycle is and set the appropriate `-timeout`, in minutes.

See the `--help` for all possible options.

[lirc]: https://github.com/hemtjanst/rodljus
