package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"math"
	"time"

	"github.com/boltdb/bolt"
)

/*
we store a temp value using actual time
*/

const dbName = "heating.db"

var db *bolt.DB

func ensureDBOpen() {
	var err error
	if db == nil {
		log.Printf("Opening DB %s\n", dbName)
		db, err = bolt.Open(dbName, 0600, nil)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func closeDb() {
	ensureDBOpen()
	db.Close()
	log.Printf("Closed DB %s\n", dbName)
}

func getOrCreateBucket(bucket string, tx *bolt.Tx) *bolt.Bucket {
	b := tx.Bucket([]byte(bucket))

	// Create Bucket on the fly
	if b == nil {
		log.Printf("Create Bucket %s ...\n", bucket)
		var err error
		b, err = tx.CreateBucket([]byte(bucket))
		if err != nil {
			log.Fatalf("Error create bucket: %s", err)
		}
	}
	return b
}

func save(bucket string, value float32) {
	ensureDBOpen()

	t := time.Now().Format(time.RFC3339)
	log.Printf("%s: %s -> %f\n", t, bucket, value)

	db.Update(func(tx *bolt.Tx) error {
		b := getOrCreateBucket(bucket, tx)
		return b.Put([]byte(t), float32bytes(value))
	})
}

func getValuesBetween(bucket, min, max string) map[string]float32 {
	var s map[string]float32
	db.View(func(tx *bolt.Tx) error {
		c := getOrCreateBucket(bucket, tx).Cursor()
		for k, v := c.Seek([]byte(min)); k != nil && bytes.Compare(k, []byte(max)) <= 0; k, v = c.Next() {
			s[string(k[:])] = float32frombytes(v)
		}
		return nil
	})
	return s
}

func countValues(bucket string) int {
	var retval int
	db.View(func(tx *bolt.Tx) error {
		b := getOrCreateBucket(bucket, tx)
		retval = b.Stats().KeyN
		log.Println("count: " + string(retval))
		return nil
	})
	return retval
}

func float32frombytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func float32bytes(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}
