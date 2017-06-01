package main

import (
	"log"
	"net/http"
	"strconv"

	"gopkg.in/gin-gonic/gin.v1"
)

func sensorInfoHandler(c *gin.Context) {
	/*
		- Anzahl der Werte
		- Minimum mit Datum
		- Maximum mit Datum
		- Minimum der letzten Woche, Monat, Jahr
		- Maximum der letzten Woche, Monat, Jahr
		- Avg
	*/
	//Ívalues := c.Query("lastvalues")
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
	/*
		Zeitraumabfrage der Werte
	*/
	// // Our time range spans the 90's decade.
	// min := []byte("1990-01-01T00:00:00Z")
	// max := []byte("2000-01-01T00:00:00Z")

}

// SensorData Temperature
type SensorData struct {
	SensorName string `json:"name"`
	Value      string `json:"value"`
}

func sensorUpdateHandler(c *gin.Context) {
	//	save("testsensor", 0.1)
	var j SensorData
	if e := c.Bind(&j); e == nil {
		value, _ := strconv.ParseFloat(j.Value, 32)
		save(j.SensorName, float32(value))
		c.JSON(http.StatusOK, gin.H{"status": j.Value})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": e})
	}
}
