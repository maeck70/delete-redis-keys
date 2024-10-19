package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

type Data_t struct {
	Data       string    `json:"data"`
	UpdateDTTM time.Time `json:"update-dttm"`
}

func main() {
	// Connect to Redis Server
	conn, err := redis.Dial("tcp", "10.0.0.180:6379", redis.DialPassword(""))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// createTestKeys(conn)
	// deleteAllKeys(conn)
	deleteExpiredKeys(conn, time.Minute)
	// printAllKeys(conn)
}

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

func deleteAllKeys(conn redis.Conn) {
	dataSet := getAllKeys(conn)

	// Loop through all keys
	for key, _ := range dataSet {
		// Delete the key
		_, err := conn.Do("DEL", key)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Deleted key: %s\n", key)
	}
}

func deleteExpiredKeys(conn redis.Conn, expiration time.Duration) {
	dataSet := getAllKeys(conn)

	cnt := 0

	// Loop through all keys
	for key, data := range dataSet {
		checkDttm := time.Now()
		updateDTTM := data.UpdateDTTM.Add(expiration)

		fmt.Printf("UpdateDTTM: %s\n", updateDTTM)
		fmt.Printf("CheckDTTM:  %s\n", checkDttm)
		fmt.Printf("Delete: %v\n", updateDTTM.Before(checkDttm))

		if updateDTTM.Before(checkDttm) {
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

func printAllKeys(conn redis.Conn) {
	dataSet := getAllKeys(conn)

	// Loop through all keys
	for key, value := range dataSet {
		fmt.Printf("Key: %s, Value: %s, UpdateDTTM: %s\n", key, value.Data, value.UpdateDTTM)
	}
}
