package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/gin-gonic/gin.v1"
)

const user = "foo"
const pass = "bar"

func handlCtrlC() {
	// Ctrl+C Handler
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nCleaning up...\n")
		closeDb()
		os.Exit(1)
	}()
}

func main() {
	handlCtrlC()

	log.Println("Listening on :8080")

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	//every call needs to be authenticated
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{user: pass}))

	// get a List of all Sensors
	authorized.GET("/sensors", sensorsHandler)

	// get info about a sensor
	authorized.GET("/sensor/:sensorname/info", sensorInfoHandler)

	//authorized.GET("/:sensorname/values", sensorValueHandler)

	// get last value of a sensor, optional wit parameter count multiple values
	authorized.GET("/sensor/:sensorname/lastvalue", sensorLastValueHandler)

	// update a Value of the sensor (with date of server)
	authorized.PUT("/sensor/:sensorname", sensorUpdateHandler)

	router.Run()
}
