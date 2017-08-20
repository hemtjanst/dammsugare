# Dammsugare

Dammsugare exposes a switch to HomeKit that allows turning a robot vacuum
on or off. By default it is configured for a Eufy RoboVac 11.

**Note**: This relies on [mqtt2lirc][lirc]. The RoboVac is fairly dumb, it
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

## Configuration

Because the robot cannot signal back to us when it's done the timeout is used
to flip the `on` state back to off. Try to guestimate how long the cleaning
cycle is and set the appropriate `-timeout`, in minutes.

See the `--help` for all possible options.

[lirc]: https://github.com/hemtjanst/mqtt2lirc
