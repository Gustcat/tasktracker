package kafka_consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Gustcat/shared-lib/kafka_common"
	"github.com/Gustcat/task-server/internal/config"
	"github.com/Gustcat/task-server/internal/service"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader      *kafka.Reader
	taskService service.TaskService
}

func NewConsumer(cfg config.ConsumerConfig, taskService service.TaskService) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: cfg.BrokerAddrs,
			Topic:   kafka_common.UserDeletedTopic,
		}),
		taskService: taskService,
	}
}

func (c *Consumer) Run(ctx context.Context) {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			fmt.Println("Consumer is running")
			//log.Println("error reading message:", err)
			continue
		}

		c.handleMessage(ctx, m)
	}
}

func (c *Consumer) handleMessage(ctx context.Context, m kafka.Message) {
	fmt.Println("handleMessage is running")
	var payload kafka_common.UserDeletedPayload
	err := json.Unmarshal(m.Value, &payload)
	if err != nil {
		fmt.Printf("error decoding: %v\n", err)
		//logger.Error()
	}

	go func() {
		fmt.Println("DeleteUserFromObservers is running")
		if err := c.taskService.DeleteUserFromObservers(ctx, payload.UserID); err != nil {
			// TODO: канал для обработки ошибок
		}
	}()

	go func() {
		if err := c.taskService.DeleteUserFromAuthors(ctx, payload.UserID); err != nil {
			// TODO: канал для обработки ошибок
		}
	}()

	go func() {
		if err := c.taskService.DeleteUserFromOperators(ctx, payload.UserID); err != nil {
			// TODO: канал для обработки ошибок
		}
	}()
}
