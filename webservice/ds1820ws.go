package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/gin-gonic/gin.v1"
)

/*
TODO: 
- Basic Auth soll bei webanwendung nicht kommen -> müsste dann eigentlich im JS sein, dass sieht aber der Client wieder  
- REST CAll wo man die values als Liste sieht
*/

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

	//Handler for every call -> needs to be authenticated
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{user: pass}))

	// Static Router with CORS enabled
	staticRouter := authorized.Group("/static", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
	})
	staticRouter.Static("/", "./static/")

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
