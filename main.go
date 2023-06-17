package main

// Go-Redis is a popular Redis client library for Go that provides a simple and easy-to-use interface for interacting with Redis from Go programs.
import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
)

func main() {

	cacheTTL := 300
	fmt.Println("Enter which set to run single or cluster redis: 0 for single , 1 for cluster : ")

	// var then variable name then variable type
	var input int8
	// Taking input from user
	fmt.Scanln(&input)

	if input == 0 {
		// Create a Redis client
		client := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})
		// Set a value in the cache
		err := client.Set("foo", "bar", time.Second*time.Duration(cacheTTL)).Err()
		if err != nil {
			log.Fatal(err)
		}
		// Retrieve a value from the cache
		val, err := client.Get("foo").Result()
		if err == redis.Nil {
			fmt.Println("key not found or TTL expired")
		} else if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("foo:", val)
		}

		// Wait for a few seconds to simulate cache expiration
		fmt.Println("Wait for Cache to expire after TTL duration :", cacheTTL, "seconds")
		time.Sleep(time.Duration(cacheTTL+1) * time.Second)

		// Retrieve the same value from the cache after expiration
		val, err = client.Get("foo").Result()
		if err == redis.Nil {
			fmt.Println("key not found")
		} else if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("foo:", val)
		}
	} else {

		// Example with REDIS CLUSTER
		// // Create a Redis CLUSTER client
		cluster := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: []string{"10.101.1.2:7000", "10.101.1.4:7002", "10.101.1.3:7001", "10.101.1.6:7004", "10.101.1.7:7005", "10.101.1.5:7003"},
		})

		err := cluster.Ping().Err()
		if err != nil {
			log.Fatal(err)
		}
		// Set a value in the cache for key=foo value=bar TTL=1min
		err = cluster.Set("foo", "barValueInredisCluster", time.Second*time.Duration(cacheTTL)).Err()
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Successfully set the value of key foo")
		}

		// Retrieve a value from the cache
		val, err := cluster.Get("foo").Result()
		if err == redis.Nil {
			fmt.Println("key not found")
		} else if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Value retrieved from cluster for key foo:", val)
		}

		// Wait for a few seconds to simulate cache expiration
		fmt.Println("Wait for Cache to expire after TTL duration :", cacheTTL, "seconds")
		time.Sleep(time.Duration(cacheTTL+1) * time.Second)

		// Retrieve the same value from the cache after expiration
		val, err = cluster.Get("foo").Result()
		if err == redis.Nil {
			fmt.Println("key not found as TTL expired")
		} else if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("foo:", val)
		}

		err = cluster.Del("key").Err()
		if err != nil {
			log.Fatal(err)
		}
	}

}
