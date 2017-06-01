package main

import (
	"fmt"
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

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	//every call needs to be authenticated
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{user: pass}))

	authorized.GET("/sensor/:sensorname/info", sensorInfoHandler)
	authorized.GET("/sensor/:sensorname/values", sensorValueHandler)
	authorized.PUT("/sensor/:sensorname", sensorUpdateHandler)

	router.Run()
}
