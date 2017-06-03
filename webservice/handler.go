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

func sensorInfoHandler(c *gin.Context) {
	/*
		- Anzahl der Werte
		- Minimum mit Datum
		- Maximum mit Datum
		- Minimum der letzten Woche, Monat, Jahr
		- Maximum der letzten Woche, Monat, Jahr
		- Avg
	*/
	name := c.Param(":sensorname")
	values := countValues(name)
	r := gin.H{"message": string(values)}
	log.Println(r)
	if values >= 0 {
		log.Println("values: " + string(values))
		c.JSON(200, r)
	} else {
		c.JSON(202, r)
	}
}

func sensorValueHandler(c *gin.Context) {
	name := c.Param("sensorname")
	lastvalue := c.Query("lastvalue")

	if lastvalue == "true" {
		date, value := getLastValues(name)
		r := gin.H{"date": date, "value": value}
		c.JSON(200, r)
	}
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
		value, _ := strconv.ParseFloat(j.Value, 32)
		save(name, float32(value))
		c.JSON(http.StatusOK, gin.H{"status": j.Value})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": e})
	}
}
