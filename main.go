package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func main() {
	// Connect to Redis Server
	conn, err := redis.Dial("tcp", "10.0.0.180:6379", redis.DialPassword(""))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Get all keys
	keys, err := redis.Strings(conn.Do("KEYS", "*"))
	if err != nil {
		panic(err)
	}

	// Print all keys
	for _, key := range keys {
		/*
			// Get the value of the key
			value, err := redis.String(conn.Do("GET", key))
			if err != nil {
				panic(err)
			}
			println(value)
		*/

		// Delete the key
		_, err = conn.Do("DEL", key)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Deleted key: %s\n", key)
	}
}
