package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
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

	url := "http://restapi3.apiary.io/notes"
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func init() {
	list = flag.Bool("list", false, "list sensors")
	transmit = flag.Bool("transmit", false, "Transmit sensors")
	flag.Parse()
}

func main() {

	if len(flag.Args()) <= 1 {
		flag.PrintDefaults()
		log.Fatal("Keine Argumente angegeben")
	}

	sensorDirs, err := ioutil.ReadDir(baseDir)
	if err != nil {
		log.Fatal(err)
	}

	var sensors []string
	for _, f := range sensorDirs {
		if m, _ := regexp.Match("\\d{2}-[0-9a-z]{12}", []byte(f.Name())); m {
			sensors = append(sensors, f.Name())
		}
	}

	fmt.Println(os.Args)

	for _, n := range sensors {
		if *list {
			fmt.Println(n)
		}
		if *transmit {
			doTransmit(n)
		}
	}
}
