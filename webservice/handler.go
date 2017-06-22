package main

import (
	"log"
	"net/http"
	"strconv"

	"gopkg.in/gin-gonic/gin.v1"
)

// SensorData Temperature
type SensorUpdateData struct {
	Value string `json:"value"`
}

func sensorsHandler(c *gin.Context) {
	log.Println("Sensors Info")
	c.JSON(200, getBuckets())
}

func sensorInfoHandler(c *gin.Context) {
	/*
		- Anzahl der Werte
		- Minimum mit Datum
		- Maximum mit Datum
		- Minimum der letzten Woche, Monat, Jahr
		- Maximum der letzten Woche, Monat, Jahr
		- Avg
	*/
	name := c.Param("sensorname")
	values := countValues(name)
	r := gin.H{"valuecount": strconv.Itoa(values)}
	log.Println(r)
	c.JSON(202, r)
}

func sensorValueHandler(c *gin.Context) {
	//name := c.Param("sensorname")
	/*
		Zeitraumabfrage der Werte
	*/
	// // Our time range spans the 90's decade.
	// min := []byte("1990-01-01T00:00:00Z")
	// max := []byte("2000-01-01T00:00:00Z")

}

func sensorUpdateHandler(c *gin.Context) {
	var j SensorUpdateData
	name := c.Param("sensorname")
	if e := c.Bind(&j); e == nil {
		log.Print("Data")
		log.Println(j)
		value, _ := strconv.ParseFloat(j.Value, 32)
		save(name, float32(value))
		c.JSON(http.StatusOK, gin.H{"status": j.Value})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": e})
	}
}

func sensorLastValueHandler(c *gin.Context) {
	name := c.Param("sensorname")
	count := c.Query("count")
	if count != "" {
		countval, _ := strconv.Atoi(count)
		c.JSON(http.StatusOK, getNLastValues(name, countval))
	} else {
		r := getNLastValues(name, 1)[0]
		c.JSON(http.StatusOK, r)
	}
}
