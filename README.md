# Cache with Golang and Redis
Go-Redis is a popular Redis client library for Go that provides a simple and easy-to-use interface for interacting with Redis from Go programs.

here we use redis latest image and run the redis container
```docker run --name my-redis -p 6379:6379 -d redis```

Go inside the container and test redis-clil
```docker exec -it my-redis sh```
```
redis-cli
127.0.0.1:6379> ping
PONG
127.0.0.1:6379> set name Monica
OK
127.0.0.1:6379> get name
"Monica"
```
Run the program
```go run main.go```
You can fetch the value from code or from redis-cli inside of the container
```redis-cli
127.0.0.1:6379> get foo
"bar"
```

# Data is persistent
```➜ docker stop my-redis
my-redis

➜ redis-cli
Could not connect to Redis at 127.0.0.1:6379: Connection refused
not connected> exit

➜ docker start my-redis
my-redis

➜ redis-cli
127.0.0.1:6379> get name
"Monica"
```
You can also Volume mount the data to host VM/Machine

# Distributed Cache with Golang and Redis

Redis Cluster is a distributed implementation of Redis that allows for horizontal scaling and high availability. It uses a sharding mechanism to partition the dataset across multiple nodes, ensuring that each node only holds a subset of the total data. Redis Cluster supports automatic node discovery, failover, and rebalancing, making it easy to scale the cluster up or down without any downtime.

1. Pull image version 6.x and with 7.x  and above the go libaray ""github.com/go-redis/redis"" is not compatible 
with redis image version 6 and above try using GOlang compatible libary https://github.com/redis/go-redis "

```docker pull redis:6.0.19-alpine3.18 ```
So, we will be creating six folders, each having a ‘redis.conf’ file, that will be used to 

2. Create docker instances.
```mkdir 7000 7001 7002 7003 7004 7005```

3. create redis.conf
Remember to keep this port same as the folder name. For e.g. In folder 7001, the port will be 7001, in folder 7002, the port will be 7002 and so on.
This <server-IP> must be replaced by the IP of the server we are using. This is to make our Redis instance accessible from the outside, where that server is also accessible.

4. Run the containers as:

```docker run -v $(pwd)/7000/redis.conf:/usr/local/etc/redis/redis.conf -d --net=red_cluster -p 7000:7000 --name myredis-0 redis:6.0.19-alpine3.18 redis-server /usr/local/etc/redis/redis.conf```

```docker run -v $(pwd)/7001/redis.conf:/usr/local/etc/redis/redis.conf -d --net=red_cluster -p 7001:7001 --name myredis-1 redis:6.0.19-alpine3.18 redis-server /usr/local/etc/redis/redis.conf```

Similalry for 7002, 7003 , 7004 and 7005 so you should have 6 containers running.
and so on .....
....

5. Now we have containers ready now , will create cluster from it
```docker exec -it <container_id> sh ```
Run the command
```redis-cli --cluster create <server-IP>:7000 <server-IP>:7001 <server-IP>:7002 <server-IP>:7003 <server-IP>:7004 <server-IP>:7005 --cluster-replicas 1```

you need the internal container's IP here use command
```docker inspect -f '{{ (index .NetworkSettings.Networks "red_cluster").IPAddress }}' myredis-0```
10.101.1.2 // for example for myredis-0 the IP is 10.101.1.2

OR with server-ip as :
```redis-cli --cluster create 10.101.1.2:7000 10.101.1.3:7001 10.101.1.4:7002 10.101.1.5:7003 10.101.1.6:7004 10.101.1.7:7005 --cluster-replicas 1```

--cluster-replicas 1 will create three nodes as master nodes and the other three as their slave nodes.

OUTPUT of Command:
>>> Performing hash slots allocation on 6 nodes...
>>> Performing Cluster Check
.......
>>> Trying to optimize slaves allocation for anti-affinity
[WARNING] Some slaves are in the same host as their master
[OK] All nodes agree about slots configuration.
>>> Check for open slots...
>>> Check slots coverage...
[OK] All 16384 slots covered.

Now you can run the main.go for REDIS CLUSTER
```go run main.go```

# data persistent in Cluster, increase the TTL to test below scenario
```docker exec -it myredis-0 sh``` // Go to master
```/data # redis-cli -p 7000```
```127.0.0.1:7000> get foo```
(error) MOVED 12182 10.101.1.4:7002 // it shows the key is in myredis-2
127.0.0.1:7000> exit
/data # exit 
```docker exec -it myredis-2 sh```
```/data # redis-cli -p 7002```
```127.0.0.1:7002> get foo```
```"barValueInredisCluster"```

Also note the default ```redis-cli``` tries to conect at localhost:6379 but we have changed the ports
redis-cli --> response is Could not connect to Redis at 127.0.0.1:6379: Connection refused
not connected> exit
redis-cli -p 7002
127.0.0.1:7002> get foo
"barValueInredisCluster"

# IMPORTANT , now will delete myredis-0 and myredis-1 both the masters to leave only 1 master
docker exec -it myredis-1 sh
/data # redis-cli -p 7001
127.0.0.1:7001> get foo
(error) MOVED 12182 10.101.1.5:7003
127.0.0.1:7001> exit
/data # exit
user1@ps-vm1:~$ docker exec -it myredis-3 sh // NOTE: Now the value is on myredis-3 ...
/data # redis-cli -p 7003
127.0.0.1:7003> get foo
"barValueInredisCluster"
127.0.0.1:7003> 
Also with 1 master and 3 workers the cluster is working fine.

# Redis Sentinel 

Redis Sentinel is a high-availability solution for Redis that provides automatic failover and monitoring of Redis instances. It works by running a set of Sentinel processes that continuously monitor the Redis instances and perform failover operations when needed. Sentinel can be used to ensure that your distributed cache is always available and can handle failures gracefully.
To use Redis Sentinel with Go-Redis, you first need to create a sentinel client that connects to the Redis Sentinel instances. You can do this by specifying the addresses of the Sentinel instances using the NewFailoverClient function, as shown below:
func main() {
    sentinel := redis.NewFailoverClient(&redis.FailoverOptions{
        SentinelAddrs: []string{"sentinel1:26379", "sentinel2:26379", "sentinel3:26379"},
        MasterName:    "mymaster",
    })

    err := sentinel.Ping().Err()
    if err != nil {
        log.Fatal(err)
    }
// Use the sentinel client to interact with the Redis master
// ...
}

Reference:
https://www.dltlabs.com/blog/how-to-setup-configure-a-redis-cluster-easily-573120
https://medium.com/@rajamanohar.mummidi/distributed-caching-in-go-fcacafafe819
# redisGolang
