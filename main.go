package main

// Go-Redis is a popular Redis client library for Go that provides a simple and easy-to-use interface for interacting with Redis from Go programs.
import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
)

func main() {

	cacheTTL := 2
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

	// Another example using CLUSTER
	// // Create a Redis CLUSTER client
	// cluster := redis.NewClusterClient(&redis.ClusterOptions{
	// 	Addrs: []string{"10.145.70.97:7000", "10.145.70.97:7001"},
	// })

	// err := cluster.Ping().Err()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // Set a value in the cache for key=foo value=bar TTL=1min
	// err = cluster.Set("foo", "bar", 5*time.Second).Err()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Retrieve a value from the cache
	// val, err := cluster.Get("foo").Result()
	// if err == redis.Nil {
	// 	fmt.Println("key not found")
	// } else if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	fmt.Println("foo:", val)
	// }

	// // Wait for a few seconds to simulate cache expiration
	// time.Sleep(7 * time.Second)

	// // Retrieve the same value from the cache after expiration
	// val, err = cluster.Get("foo").Result()
	// if err == redis.Nil {
	// 	fmt.Println("key not found")
	// } else if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	fmt.Println("foo:", val)
	// }

	// err = cluster.Del("key").Err()
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
