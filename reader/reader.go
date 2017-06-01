package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

/*
TODO:
- Commandline Parser
	- list all Sensors
	- Transmit all Sensors
*/

const baseDir = "/sys/bus/w1/devices"

var list, transmit *bool

func doTransmit(name string) {
	fmt.Printf("We Transmit value of %s\n", name)
}

func init() {
	list = flag.Bool("list", false, "list sensors")
	transmit = flag.Bool("transmit", false, "Transmit sensors")
	flag.Parse()
}

func main() {

	sensors, err := ioutil.ReadDir(baseDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range sensors {
		if *list {
			fmt.Println(f.Name())
		}
		if *transmit {
			doTransmit(f.Name())
		}
	}
}
