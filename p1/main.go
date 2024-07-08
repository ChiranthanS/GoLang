package main

import (
    "context"
    "github.com/confluentinc/confluent-kafka-go/kafka"
    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"
    "log"
    "net/http"
    "os"
    "strings"
)

var (
    redisClient *redis.Client
    kafkaBroker = os.Getenv("KAFKA_BROKER")
    redisAddr   = os.Getenv("REDIS_ADDR")
    redisPass   = os.Getenv("REDIS_PASS")
    ctx         = context.Background()
)

func main() {
    // Initialize Redis client
    redisClient = redis.NewClient(&redis.Options{
        Addr:     redisAddr,
        Password: redisPass, // no password set
        DB:       0,         // use default DB
    })

    // Check Redis connection
    _, err := redisClient.Ping(ctx).Result()
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
    log.Println("Connected to Redis")

    r := gin.Default()

    r.POST("/data", postDataHandler)
    r.GET("/data", getDataHandler)

    go kafkaConsumer()

    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}

func postDataHandler(c *gin.Context) {
    var jsonData map[string]string
    if err := c.ShouldBindJSON(&jsonData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }

    for key, value := range jsonData {
        go kafkaProducer(key, value)
    }

    c.JSON(http.StatusOK, gin.H{"status": "data received"})
}

func getDataHandler(c *gin.Context) {
    key := c.Query("key")

    value, err := redisClient.Get(ctx, key).Result()
    if err == redis.Nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
        return
    } else if err != nil {
        log.Printf("Failed to get key from redis: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get key from redis"})
        return
    }

    go kafkaProducer("get", key)

    c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
}

func kafkaProducer(key, value string) {
    p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaBroker})
    if err != nil {
        log.Printf("Failed to create producer: %v", err)
        return
    }
    defer p.Close()

    topic := "myTopic"
    data := key + ":" + value

    p.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
        Value:          []byte(data),
    }, nil)
}

func kafkaConsumer() {
    c, err := kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers": kafkaBroker,
        "group.id":          "myGroup",
        "auto.offset.reset": "earliest",
    })

    if err != nil {
        log.Printf("Failed to create consumer: %v", err)
        return
    }
    defer c.Close()

    c.SubscribeTopics([]string{"myTopic"}, nil)

    for {
        msg, err := c.ReadMessage(-1)
        if err == nil {
            go handleKafkaMessage(msg)
        } else {
            log.Printf("Consumer error: %v (%v)\n", err, msg)
        }
    }
}

func handleKafkaMessage(msg *kafka.Message) {
    data := string(msg.Value)
    parts := strings.SplitN(data, ":", 2)

    if len(parts) != 2 {
        log.Printf("Invalid message format: %s", data)
        return
    }

    key, value := parts[0], parts[1]

    if key == "get" {
        log.Printf("GET request received for key: %s", value)
        // Handle GET request Kafka message if needed
    } else {
        err := redisClient.Set(ctx, key, value, 0).Err()
        if err != nil {
            log.Printf("Failed to store key-value in redis: %v", err)
            return
        }
        log.Printf("Stored key-value in redis: %s=%s", key, value)
    }
}
