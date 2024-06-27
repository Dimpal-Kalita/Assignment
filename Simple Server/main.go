package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/IBM/sarama"
)

var (
	rdb           *redis.Client
	ctx           = context.Background()
	kafkaProducer sarama.SyncProducer
	kafkaConsumer sarama.Consumer
)

const (
	redisAddr  = "localhost:6379"
	kafkaAddr  = "localhost:29092"
	kafkaTopic = "my_topic"
)

func main() {
	// Initialize Redis client
	rdb = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Initialize Kafka producer
	producer, err := sarama.NewSyncProducer([]string{kafkaAddr}, nil)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}
	kafkaProducer = producer
	defer kafkaProducer.Close()

	// Initialize Kafka consumer
	consumer, err := sarama.NewConsumer([]string{kafkaAddr}, nil)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}
	kafkaConsumer = consumer
	defer kafkaConsumer.Close()

	// Set up Gin router
	router := gin.Default()

	router.POST("/push", pushHandler)
	router.GET("/get", getHandler)

	// Start background consumer
	go consumeFromKafka()

	// Start server
	router.Run(":8080")
}

func pushHandler(c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.BindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go produceToKafka(jsonData)

	c.JSON(http.StatusOK, gin.H{"status": "data pushed to Kafka"})
}

func getHandler(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key query parameter is required"})
		return
	}

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"key": key, "value": val})
}

func produceToKafka(data map[string]interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling data to JSON: %v", err)
		return
	}
	msg := &sarama.ProducerMessage{
		Topic: kafkaTopic,
		Value: sarama.StringEncoder(jsonData),
	}
	_, _, err = kafkaProducer.SendMessage(msg)
	if err != nil {
		log.Printf("Error producing message to Kafka: %v", err)
	}
}

func consumeFromKafka() {
	partitionConsumer, err := kafkaConsumer.ConsumePartition(kafkaTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error starting Kafka consumer: %v", err)
	}
	defer partitionConsumer.Close()

	for message := range partitionConsumer.Messages() {
		go processKafkaMessage(message)
	}
}

func processKafkaMessage(message *sarama.ConsumerMessage) {

	// Parse the JSON message
	var data map[string]interface{}
	if err := json.Unmarshal(message.Value, &data); err != nil {
		log.Printf("Error unmarshaling JSON message: %v", err)
		return
	}

	// Extract key and value
	var key, value string
	var keyOk, valueOk bool

	if k, ok := data["key"].(string); ok {
		key = k
		keyOk = true
	} else if k, ok := data["key"].(float64); ok {
		key = strconv.FormatFloat(k, 'f', -1, 64)
		keyOk = true
	}

	if v, ok := data["value"].(string); ok {
		value = v
		valueOk = true
	} else if v, ok := data["value"].(float64); ok {
		value = strconv.FormatFloat(v, 'f', -1, 64)
		valueOk = true
	}

	if !keyOk || !valueOk {
		// Additional debugging information
		log.Printf("Key or value not found or not a string/number: keyOk=%v valueOk=%v data=%v", keyOk, valueOk, data)
		return
	}

	err := rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		log.Printf("Error storing key-value in Redis: %v", err)
	}

	// Store the latest message
	latestMessage, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling latest message to JSON: %v", err)
		return
	}

	err = rdb.Set(ctx, "latest_message", latestMessage, 0).Err()
	if err != nil {
		log.Printf("Error storing latest message in Redis: %v", err)
	}
}
