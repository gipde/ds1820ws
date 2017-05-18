package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"gopkg.in/gin-gonic/gin.v1"
)

func initdb() *bolt.DB {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("MyBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	return db
}

var i = 0

func save(db *bolt.DB) {
	log.Println("We save a value")
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		err := b.Put([]byte(strconv.Itoa(1000+i)), []byte(strconv.Itoa(i)))
		return err
	})
}

func testDB(db *bolt.DB) {
	i++
	log.Printf("--MARK-- %d", i)
	time.Sleep(1000 * time.Millisecond)

	save(db)

	db.View(func(tx *bolt.Tx) error {
		log.Println("we view into the db")
		b := tx.Bucket([]byte("MyBucket"))

		b.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})
		return nil
	})
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pond j",
		})
	})
	go func() {
		var db = initdb()
		for {
			testDB(db)
		}
	}()

	r.Run() // listen and serve on 0.0.0.0:8080
}
