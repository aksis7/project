package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"wb/internal/cache"
	"wb/internal/order"

	"github.com/cenkalti/backoff/v4"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

func ConsumeKafkaMessages(ctx context.Context, db *gorm.DB, rdb *cache.RedisClient) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"wb-kafka-1:9092"},
		Topic:     "orders",
		Partition: 0,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			logger.Error("could not read message from Kafka", slog.String("error", err.Error()))
			continue
		}

		if len(m.Value) == 0 {
			logger.Warn("received empty message from Kafka")
			continue
		}

		var orderData order.Order
		if err := json.Unmarshal(m.Value, &orderData); err != nil {
			logger.Error("error unmarshalling kafka message",
				slog.String("error", err.Error()), slog.String("message", string(m.Value)))
			continue
		}

		// Используем функцию валидации из validate.go
		if err := order.ValidateOrder(&orderData, logger); err != nil {
			logger.Error("invalid order data received from Kafka", slog.String("error", err.Error()))
			continue
		}

		// Дополнительная проверка дубликатов в кэше и базе данных
		_, err = rdb.Client.Get(ctx, orderData.OrderUID).Result()
		if err == nil {
			logger.Warn("duplicate order detected in Redis", slog.String("order_uid", orderData.OrderUID))
			continue
		}

		var existingOrder order.Order
		if err := db.First(&existingOrder, "order_uid = ?", orderData.OrderUID).Error; err == nil {
			logger.Warn("duplicate order detected in database", slog.String("order_uid", orderData.OrderUID))
			continue
		}

		saveOrderWithRetry(orderData, db)

		// Сохраняем заказ в кэш и Redis
		rdb.Client.Set(ctx, orderData.OrderUID, m.Value, 0)

		logger.Info("processed order", slog.String("order_uid", orderData.OrderUID))
	}
}

func saveOrderWithRetry(order order.Order, db *gorm.DB) {
	operation := func() error {
		return db.Save(&order).Error
	}
	logger.Info("Order Data Before Save", slog.Any("order", order))
	expBackoff := backoff.NewExponentialBackOff()
	if err := backoff.Retry(operation, expBackoff); err != nil {
		logger.Error("failed to save order after retries",
			slog.String("order_uid", order.OrderUID), slog.String("error", err.Error()))
	} else {
		logger.Info("order saved successfully", slog.String("order_uid", order.OrderUID))
	}
}
