package main

import (
	"fmt"
	"runtime"
	"schneidernet/ds1820ws/adder"

	"gopkg.in/gin-gonic/gin.v1"
)

func startGin() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/ping", pingHandler)
	fmt.Println(adder.Add(1, 2))
	router.Run(":8080")
}

func pingHandler(c *gin.Context) {
	c.Query("hallo")
}

func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)

}

func main() {
	ConfigRuntime()
	startGin()

	// r := gin.Default()
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })

	// // Get user value
	// r.GET("/user/:name", func(c *gin.Context) {
	// 	user := c.Params.ByName("name")
	// 	value, ok := DB[user]
	// 	if ok {
	// 		c.JSON(200, gin.H{"user": user, "value": value})
	// 	} else {
	// 		c.JSON(200, gin.H{"user": user, "status": "no value"})
	// 	}
	// })

	// authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
	// 	"foo":  "bar", // user:foo password:bar
	// 	"manu": "123", // user:manu password:123
	// }))

	// authorized.POST("admin", func(c *gin.Context) {
	// 	user := c.MustGet(gin.AuthUserKey).(string)

	// 	// Parse JSON
	// 	var json struct {
	// 		Value string `json:"value" binding:"required"`
	// 	}

	// 	if c.Bind(&json) == nil {
	// 		DB[user] = json.Value
	// 		c.JSON(200, gin.H{"status": "ok"})
	// 	}
	// })

	// // // DB Goroutine
	// // go func() {
	// // 	var db = initdb()
	// // 	for {
	// // 		testDB(db)
	// // 	}
	// // }()

	// r.Run() // listen and serve on 0.0.0.0:8080
}
