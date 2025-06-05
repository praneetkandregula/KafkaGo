package main

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// places orders, even if order processor is down
type OrderPlacer struct {
	producer  *kafka.Producer
	topic     string
	deliverCh chan kafka.Event
}

func NewOrderPlacer(p *kafka.Producer, topic string) *OrderPlacer {
	return &OrderPlacer{
		producer:  p,
		topic:     topic,
		deliverCh: make(chan kafka.Event, 10000),
	}
}

func (op *OrderPlacer) placeOrder(orderType string, size int) error {
	var (
		format  = fmt.Sprintf("%s - %d", orderType, size)
		payload = []byte(format)
	)

	err := op.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &op.topic,
			Partition: kafka.PartitionAny,
		},
		Value: payload,
	},
		op.deliverCh,
	)
	if err != nil {
		log.Fatal(err)
	}

	<-op.deliverCh
	fmt.Printf("placed order on queue %s\n", format)
	return nil
}

func main() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         "foo",
		"acks":              "all",
	})
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
	}

	op := NewOrderPlacer(p, "praneet")
	for i := 0; i < 1000; i++ {
		if err := op.placeOrder("newOrders!", i+1); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second * 3)
	}
}
