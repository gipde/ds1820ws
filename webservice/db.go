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

func init() {
	ensureDBOpen()
}

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
	t := time.Now().Format(time.RFC3339)

	db.Update(func(tx *bolt.Tx) error {
		b := getOrCreateBucket(bucket, tx)
		return b.Put([]byte(t), float32bytes(value))
	})
}

/*
Entry Transfer Object for Temperature Sensor Data
*/
type Entry struct {
	Date  string
	Value float32
}

func getBuckets() []string {
	var buckets []string
	db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			buckets = append(buckets, string(name))
			log.Println(string(name))
			return nil
		})
	})
	return buckets
}

func printAllValues(bucket string) {
	log.Println("All Values --BEGIN--")
	db.View(func(tx *bolt.Tx) error {
		b := getOrCreateBucket(bucket, tx)
		log.Println("Entries: ", countValues(bucket))
		b.ForEach(func(k []byte, v []byte) error {
			log.Println(string(k), ":", float32frombytes(v))
			return nil
		})
		return nil
	})
	log.Println("All Values --END--")
}

func getNLastValues(bucket string, count int) []Entry {
	var elements []Entry
	var value, date []byte

	entries := countValues(bucket)
	if count > entries {
		count = entries
	}

	db.View(func(tx *bolt.Tx) error {
		c := getOrCreateBucket(bucket, tx).Cursor()
		date, value = c.Last()
		for count > 0 {
			elements = append(elements, Entry{string(date), float32frombytes(value)})
			if count > 1 {
				date, value = c.Prev()
			}
			count--
		}
		return nil
	})
	//reversing slice
	for i, j := 0, len(elements)-1; i < j; i, j = i+1, j-1 {
		elements[i], elements[j] = elements[j], elements[i]
	}
	return elements
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
		return nil
	})
	return retval
}

func float32frombytes(bytes []byte) float32 {
	if len(bytes) == 0 {
		return 0.0
	}
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
