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


## import data with wget
```
http_proxy="http://localhost:8080" wget -qO- --user=foo 
  --password=bar --method=PUT --header "Content-Type: application/json" 
  --body-data="{\"Value\":\"1.2\"}" "http://localhost:3000/sensor/10-0008019462b6"
```

## use Rest API
```
http_proxy="http://localhost:8080" wget -qO- --user=foo --password=bar 
  --method=GET "http://localhost:3000/sensor/10-0008019462b6/values?lastvalue=true"
```
```
http_proxy="http://localhost:8080" wget -qO- --user=foo --password=bar 
  --method=GET "http://localhost:3000/sensor/10-0008019462b6/lastvalue?count=2"
```



# Unterschied zwischen sensor 1 und sensor 2
SELECT mean("T1")-mean("T2") as "Unterschied" from (SELECT "temp" as T1 FROM "heating"."autogen"."heating" where "sensor"='10-0008019462b6' ),(SELECT "temp" as T2 FROM "heating"."autogen"."heating" where "sensor"='10-000800355f27' ) WHERE time > now() - 10h  GROUP BY time(1m)

SELECT mean("T1")-mean("T2") as "Unterschied" from (SELECT "temp" as T1 FROM "heating"."autogen"."heating" where "sensor"='10-0008019462b6' ),(SELECT "temp" as T2 FROM "heating"."autogen"."heating" where "sensor"='10-000800355f27' ) WHERE time > :dashboardTime:  GROUP BY time(1m)


SELECT mean("temp") AS "mean_temp" FROM "heating"."autogen"."heating" WHERE time > :dashboardTime: AND "sensor"='10-0008019462b6' GROUP BY :interval:, "sensor"


8019462b6 Umschaltventiel Zulauf
800355f27 Umschaltventiel Heizung
8019481df Umschaltventiel Wasser
80194662f Auslauf Warmwasser
8019453e2 Warmwasser

SELECT mean("T1")-mean("T2") as "Unterschied" from (SELECT "temp" as T1 FROM "heating"."autogen"."heating" where "sensor"='10-00080194662f' ),(SELECT "temp" as T2 FROM "heating"."autogen"."heating" where "sensor"='10-0008019453e2' ) WHERE time > :dashboardTime:  GROUP BY time(1m)

SELECT mean("T1")-mean("T2") as "Unterschied" from (SELECT "temp" as T1 FROM "heating"."autogen"."heating" where "sensor"='10-0008019462b6' ),(SELECT "temp" as T2 FROM "heating"."autogen"."heating" where "sensor"='10-0008019481df' ) WHERE time > :dashboardTime:  GROUP BY time(1m)

SELECT top("temp",10) AS "mean_temp" FROM "heating"."autogen"."heating" WHERE time > :dashboardTime: AND "sensor"='10-0008019481df' 
SELECT spread(*) AS "spread_temp" FROM "heating"."autogen"."heating" WHERE time > :dashboardTime: AND "sensor"='10-0008019481df' GROUP BY time(2m)
SELECT integral("temp")  FROM "heating"."autogen"."heating" WHERE time > :dashboardTime: AND "sensor"='10-000800355f27' group by :interval:
select non_negative_derivative("temp") FROM "heating"."autogen"."heating" WHERE time > :dashboardTime: AND "sensor"='10-0008019481df'
select derivative(first) from (SELECT derivative("temp") as first FROM "heating"."autogen"."heating" WHERE time > :dashboardTime: AND "sensor"='10-0008019481df' )