package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

type Data_t struct {
	Data       string    `json:"data"`
	UpdateDTTM time.Time `json:"update-dttm"`
}

var (
	host          string
	port          int
	expiration    int
	deleteAll     bool
	deleteExpired bool
	createTest    bool
	outputKeys    bool
)

func main() {
	// Process flags
	flag.StringVar(&host, "h", "localhost", "Redis server hostname")
	flag.IntVar(&port, "p", 6379, "Redis port")
	flag.IntVar(&expiration, "e", 15, "Expiration in minutes")
	flag.BoolVar(&deleteAll, "da", false, "Delete all keys")
	flag.BoolVar(&deleteExpired, "de", false, "Delete expired keys")
	flag.BoolVar(&createTest, "c", false, "Create test keys")
	flag.BoolVar(&outputKeys, "o", false, "Output all keys")
	flag.Parse()

	// Connect to Redis Server
	conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port), redis.DialPassword(""))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	if createTest {
		createTestKeys(conn)
	}

	if deleteAll {
		deleteAllKeys(conn)
	}

	if deleteExpired {
		deleteExpiredKeys(conn, time.Duration(expiration)*time.Minute)
	}

	if outputKeys {
		printAllKeys(conn)
	}
}

// Get all keys and their values
func getAllKeys(conn redis.Conn) map[string]Data_t {
	dataSet := make(map[string]Data_t)

	// Get all keys
	keys, err := redis.Strings(conn.Do("KEYS", "*"))
	if err != nil {
		panic(err)
	}

	// Loop through all keys
	for _, key := range keys {
		// Get the value of the key
		value, err := redis.String(conn.Do("GET", key))
		if err != nil {
			panic(err)
		}

		d := Data_t{}
		err = json.Unmarshal([]byte(value), &d)
		if err != nil {
			panic(err)
		}
		dataSet[key] = d
	}

	return dataSet
}

// Delete all keys, regardless of expiration
func deleteAllKeys(conn redis.Conn) {
	dataSet := getAllKeys(conn)

	// Loop through all keys
	for key := range dataSet {
		// Delete the key
		_, err := conn.Do("DEL", key)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Deleted key: %s\n", key)
	}
}

// Delete expired keys. Provide expiration duration that will be added to the UpdateDTTM before comparing to the current time.
func deleteExpiredKeys(conn redis.Conn, expiration time.Duration) {
	dataSet := getAllKeys(conn)

	cnt := 0

	// Loop through all keys
	for key, data := range dataSet {
		checkDttm := time.Now()
		updateDTTM := data.UpdateDTTM.Add(expiration)

		if updateDTTM.Before(checkDttm) {
			fmt.Printf("Key %s is expired\n", key)
			// Delete the key
			_, err := conn.Do("DEL", key)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Deleted key: %s\n", key)
			cnt += 1
		}
	}
	fmt.Printf("Deleted %d keys\n\n", cnt)
}

// Create test keys
func createTestKeys(conn redis.Conn) {
	for i := 0; i < 10; i++ {
		// create data
		data := Data_t{
			Data:       fmt.Sprintf("value%d", i),
			UpdateDTTM: time.Now().Add(time.Duration(i*10) * time.Second),
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}

		key, err := uuid.NewUUID()
		if err != nil {
			panic(err)
		}

		_, err = conn.Do("SET", key.String(), jsonData)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Created key: %s\n\n", key.String())
	}
}

// Print all keys
func printAllKeys(conn redis.Conn) {
	dataSet := getAllKeys(conn)

	// Loop through all keys
	for key, value := range dataSet {
		fmt.Printf("Key: %s, Value: %s, UpdateDTTM: %s\n", key, value.Data, value.UpdateDTTM)
	}
}
