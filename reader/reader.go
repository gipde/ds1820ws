package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

const (
	baseDir = "/sys/bus/w1/devices"
	mydb    = "heating"
)

var list, transmit *bool
var hostname *string
var port *int
var user *string
var password *string
var delay *int

func init() {
	list = flag.Bool("list", false, "list sensors")
	transmit = flag.Bool("transmit", false, "Transmit sensors")
	hostname = flag.String("host", "2e1512f0-d590-4eed-bb41-9ad3abd03edf.pub.cloud.scaleway.com", "hostname")
	port = flag.Int("port", 8086, "Port")
	user = flag.String("user", "secret", "username")
	password = flag.String("password", "secret", "password")
	delay = flag.Int("delay", 0, "Endlosschleife Wartezeit")
	flag.Usage = usage
	flag.Parse()
}

func usage() {
	fmt.Fprintf(os.Stderr, "keine Argumente:\nUsage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func readSensorFile(f string) (float64, error) {
	file, err := os.Open(baseDir + "/" + f + "/w1_slave")
	if err != nil {
		return 0, fmt.Errorf("unable to read sensor:  %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan() // 1. Line
	if v, _ := regexp.MatchString(".*YES", scanner.Text()); !v {
		return 0, fmt.Errorf("CRC failed")
	}
	scanner.Scan() // 2. Line
	re := regexp.MustCompile(".*t=(\\d*)")
	matches := re.FindStringSubmatch(scanner.Text())

	var temp float64
	if len(matches) == 2 {
		tempInt, _ := strconv.Atoi(matches[1])
		temp = float64(tempInt) / 1000

		// if tmp==85.00 usually a read error
		if temp != 85.00 && temp != 0.00 {
			log.Printf("sensor %s: %03f", f, temp)
			return temp, nil
		}
	}
	return 0, fmt.Errorf("invalid Value at %s %0.2f", f, temp)
}

func addBatchPoint(bp client.BatchPoints, name string) error {

	errcount := 3

	var err error
	var temp float64
	for errcount > 0 {
		temp, err = readSensorFile(name)
		if err == nil {
			break
		}
		if errcount > 0 {
			log.Printf("error reading %s, next try (%s)", name, err)
		}
		errcount--
	}
	if errcount == 0 {
		return fmt.Errorf("we give up, to many errors in reading %s: %s", name, err)
	}

	tags := map[string]string{"sensor": name}
	fields := map[string]interface{}{"temp": temp}

	pt, err := client.NewPoint("heating", tags, fields, time.Now())
	if err != nil {
		return err
	}
	bp.AddPoint(pt)

	return nil

}

func transmitSensorData(sensors *[]string) error {
	// create influx client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     fmt.Sprintf("http://%s:%d", *hostname, *port),
		Username: *user,
		Password: *password,
	})

	if err != nil {
		return err
	}

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  mydb,
		Precision: "s",
	})
	if err != nil {
		return err
	}

	// add points
	sensorValues := 0
	for _, n := range *sensors {
		if err := addBatchPoint(bp, n); err == nil {
			sensorValues++
		} else {
			log.Printf("Error: %s", err)
		}
	}

	// Write the batch, if get min 1 value
	if sensorValues > 0 {
		err := c.Write(bp)
		//		log.Printf("Transmitted data")
		err2 := c.Close()
		if err2 != nil {
			return err
		}
		//		log.Printf("Connection closed")
		if err != nil {
			return err
		}
	}

	return nil

}

func main() {

	if len(os.Args) <= 1 {
		usage()
	}

	log.Println("Starting...")

	if *list {
		detectSensors()
		os.Exit(0)
	}

	// transfer loop
	for {
		err := transmitSensorData(detectSensors())
		if err != nil {
			log.Printf(err.Error())
		}

		if (*delay) < 1 {
			os.Exit(0)
		}
		time.Sleep(time.Second * time.Duration(*delay))
	}

}

func listSensors(sensors *[]string) {
	fmt.Println("we have detected the following sensors:")
	for _, n := range *sensors {
		fmt.Println(n)
	}
}

func detectSensors() *[]string {
	sensorDirs, err := ioutil.ReadDir(baseDir)
	if err != nil {
		log.Fatal("Kann nicht vom Directory lesen: " + err.Error())
	}

	var sensors []string
	for _, f := range sensorDirs {
		if m, _ := regexp.Match("10-[0-9a-z]{12}", []byte(f.Name())); m {
			sensors = append(sensors, f.Name())
		}
	}
	listSensors(&sensors)
	return &sensors
}
