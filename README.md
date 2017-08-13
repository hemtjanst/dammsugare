# Dammsugare

Dammsugare exposes a switch to HomeKit that allows turning a Eufy RoboVac 11 on
or off.

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
