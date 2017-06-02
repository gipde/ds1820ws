package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

const baseDir = "/sys/bus/w1/devices"
const url = "http://localhost/sensor/"
const user = "foo"
const pass = "bar"

var list, transmit *bool

// SensorData Temperature
type SensorData struct {
	SensorName string `json:"name"`
	Value      string `json:"value"`
}

func init() {
	list = flag.Bool("list", false, "list sensors")
	transmit = flag.Bool("transmit", false, "Transmit sensors")
	flag.Usage = usage
	flag.Parse()
}

func readSensorFile(f string) string {
	file, err := os.Open(baseDir + "/" + f + "/w1_slave")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan() // 1. Line
	if v, _ := regexp.MatchString(".*YES", scanner.Text()); !v {
		log.Fatal("CRC failed")
	}
	scanner.Scan() // 2. Line
	re := regexp.MustCompile(".*t=(\\d*)")
	matches := re.FindStringSubmatch(scanner.Text())

	var retval = "INVALID"

	if len(matches) == 2 {
		tempInt, _ := strconv.Atoi(matches[1])
		retval = strconv.FormatFloat(float64(tempInt)/1000, 'f', 2, 32)
	}
	return retval

}

func doTransmit(name string) {
	jsonStr, _ := json.Marshal(SensorData{name, readSensorFile(name)})
	fmt.Printf("Transmitting: %s\n", jsonStr)
	req, _ := http.NewRequest("POST", url+"/"+name, bytes.NewBuffer(jsonStr))
	req.SetBasicAuth(user, pass)
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

func usage() {
	fmt.Fprintf(os.Stderr, "keine Argumente:\nUsage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {

	//	fmt.Println(len(os.Args))
	if len(os.Args) <= 1 {
		usage()
	}

	sensorDirs, err := ioutil.ReadDir(baseDir)
	if err != nil {
		log.Fatal(err)
	}

	var sensors []string
	for _, f := range sensorDirs {
		if m, _ := regexp.Match("10-[0-9a-z]{12}", []byte(f.Name())); m {
			sensors = append(sensors, f.Name())
		}
	}

	for _, n := range sensors {
		if *list {
			fmt.Println(n)
		}
		if *transmit {
			doTransmit(n)
		}
	}
}
