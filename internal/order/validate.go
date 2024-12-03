package order

import (
	"fmt"
	"log/slog"
	"os"
)

// ValidateOrder - функция для проверки правильности заполнения данных заказа
func ValidateOrder(order *Order, logger *slog.Logger) error {
	if logger == nil {
		// Если логгер не передан, создаем логгер по умолчанию
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	}

	// Проверяем, указан ли уникальный идентификатор заказа (OrderUID)
	if order.OrderUID == "" {
		err := fmt.Errorf("order_uid обязателен для заполнения")
		logger.Error("Ошибка валидации", slog.String("поле", "order_uid"), slog.String("ошибка", err.Error()))
		return err
	}

	// Проверяем, указан ли трек-номер
	if order.TrackNumber == "" {
		err := fmt.Errorf("track_number обязателен для заполнения")
		logger.Error("Ошибка валидации", slog.String("поле", "track_number"), slog.String("ошибка", err.Error()))
		return err
	}

	// Проверяем, указано ли имя получателя
	if order.Delivery.Name == "" {
		err := fmt.Errorf("delivery name обязателен для заполнения")
		logger.Error("Ошибка валидации", slog.String("поле", "delivery_name"), slog.String("ошибка", err.Error()))
		return err
	}

	// Проверяем, указан ли телефон получателя
	if order.Delivery.Phone == "" {
		err := fmt.Errorf("delivery phone обязателен для заполнения")
		logger.Error("Ошибка валидации", slog.String("поле", "delivery_phone"), slog.String("ошибка", err.Error()))
		return err
	}

	// Проверяем, указана ли информация о платеже
	if order.Payment.Transaction == "" {
		err := fmt.Errorf("payment transaction обязателен для заполнения")
		logger.Error("Ошибка валидации", slog.String("поле", "payment_transaction"), slog.String("ошибка", err.Error()))
		return err
	}

	// Вы можете добавить дополнительные проверки для других полей при необходимости

	// Если все проверки пройдены, возвращаем nil (ошибок нет)
	return nil
}
