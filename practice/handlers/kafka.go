package handlers

import (
	"context"
	"fmt"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/segmentio/kafka-go"
	util "tony123.tw/util"
)

// func ConsumeMessage(nodeName string) ([]kafka.Message, error) {
// 	// get the message from kafka broker
// 	// hard-coded only process ten messages at a time
// 	bootstrapServers := "my-cluster-kafka-bootstrap.kafka.svc.cluster.local:9092"
// 	topic := "my-topic"
// 	groupID := "my-group"

// 	// Create a Kafka consumer (reader)
// 	reader := kafka.NewReader(kafka.ReaderConfig{
// 		Brokers: []string{bootstrapServers},
// 		Topic:   topic,
// 		GroupID: groupID,
// 	})

// 	defer reader.Close()
// 	// Consume messages from the topic
// 	now := time.Now()
// 	messageList := []kafka.Message{}
// 	for {
// 		fmt.Println("ready to fetch message",nodeName)
// 		msg, err := reader.FetchMessage(context.Background())
// 		fmt.Println("msg.Key: ", string(msg.Key),"msg.Value:", string(msg.Value), "nodeName: ", nodeName)
// 		if err != nil {
// 			log.Fatalf("Failed to fetch message: %v", err)
// 		}
// 		if string(msg.Key) == nodeName {
// 			// Commit the offset to acknowledge the message has been processed
// 			if err := reader.CommitMessages(context.Background(), msg); err != nil {
// 				log.Fatalf("Failed to commit message: %v", err)
// 			}
// 			messageList = append(messageList, msg)
// 		}else{
// 			break
// 		}
// 		if time.Since(now) > 30 * time.Millisecond {
// 			break
// 		}
// 	}
// 	return messageList, nil

// }
// ConsumeMessage consumes messages from the Kafka topic for the specified nodeName.
func ConsumeMessage(nodeName string) ([]kafka.Message, error) {
	bootstrapServers := "my-cluster-kafka-bootstrap.kafka.svc.cluster.local:9092"
	topic := "my-topic"
	groupID := "my-group"

	// Create a Kafka consumer (reader)
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{bootstrapServers},
		Topic:   topic,
		GroupID: groupID,
	})

	defer reader.Close()
	messageList := []kafka.Message{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("ready to fetch message", nodeName)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("context done")
			return messageList, ctx.Err()
		default:
			msg, err := reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() == context.DeadlineExceeded {
					return messageList, nil
				}
				continue
			}

			fmt.Println("msg.Key: ", string(msg.Key), "msg.Value:", string(msg.Value), "nodeName: ", nodeName)

			if string(msg.Key) == nodeName {
				if err := reader.CommitMessages(ctx, msg); err != nil {
					fmt.Printf("Failed to commit message: %v", err)
					return messageList, err
				}
				messageList = append(messageList, msg)
				fmt.Println("messageList: ", messageList)
			} 
		}
	}
}
func ProduceMessage(key string, value string) error {
	bootstrapServers := "my-cluster-kafka-bootstrap.kafka.svc.cluster.local:9092"
	topic := "my-topic"
	value = util.ModifyCheckpointToImageName(value)

	// Create a Kafka producer
	writer := kafka.Writer{
		Addr:     kafka.TCP(bootstrapServers),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		Async:    false,
	}

	// Prepare the message,key and value are read from environment variables
	message := kafka.Message{
		Key:   []byte(key), // Optional: specify a key for the message
		Value: []byte(value),
	}

	// Send the message
	err := writer.WriteMessages(context.Background(), message)
	if err != nil {
		return err
	}
	fmt.Println("message sent", key, value)

	// Close the producer
	if err := writer.Close(); err != nil {
		log.Log.Error(err, "Failed to close writer")
	}
	return nil
}
