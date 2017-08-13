# Dammsugare

Dammsugare exposes a switch to HomeKit that allows turning a Eufy RoboVac 11 on
or off.

**Note**: This relies on a soon to be open sourced component that exposes an
[LIRCd](http://www.lirc.org/html/lircd.html) to MQTT. Each topic it exposes is
named after the associated key you want to press and you publish a duration to
the topic to initiate a key press. This essentially gives you an MQTT powered
IR remote that can be used to control any device within range.

## Installation

The installation is pretty simple, `go install` it and run it. Potentially
adjust the timeout.

```
go install github.com/hemtjanst/dammsugare
dammsugare -timeout 100 -mqtt.address broker.mydomain.tld:1883
```

## Configuration

Because the robot cannot signal back to us when it's done the timeout is used
to flip the `on` state back to off. Try to guestimate how long the cleaning
cycle is and set the appropriate `-timeout`, in minutes.

See the `--help` for all possible options.
