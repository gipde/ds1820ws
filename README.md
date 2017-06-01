# ds1820ws

this program is used to store from ds1820 sensors into a embedded db
and make this accessible by a rest service

relevant path: 
```shell
/sys/bus/w1/devices/w1_bus_master1
````

## Crosscompile RaspberryPi:

```sh
env GOOS=linux GOARCH=arm GOARM=6 go build
```

## Crosscompile Linux:
```sh
env GOOS=linux GOARCH=amd64 go build
```

## Design Guides

Wäremepumpe ----- Multiple Sensors    ||   Raspberry-Pi    ||     Cloud-Service  ||    Frontend
                                            go Deamon               go Deamon             native APP
                                            REST Push               REST Pull             REST Pull
                                                                    Bolt DB               QT Creator
                                                                                          HTML5

