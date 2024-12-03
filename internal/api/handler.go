package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
	"wb/internal/cache"
	"wb/internal/order"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// APIHandler представляет структуру обработчика API
type APIHandler struct {
	DB     *gorm.DB
	Cache  *sync.Map
	Redis  *cache.RedisClient
	Logger *slog.Logger
}

// SyncCacheWithDB выполняет синхронизацию кэша с базой данных
func (h *APIHandler) SyncCacheWithDB() {
	go func() {
		for {
			// Загружаем все заказы из базы данных
			var orders []order.Order
			err := h.DB.
				Preload("Items").
				Preload("Delivery").
				Preload("Payment").
				Find(&orders).Error

			if err != nil {
				// Логируем ошибку и ждем перед следующим повтором
				h.Logger.Error("Ошибка синхронизации кэша с базой данных",
					slog.String("error", err.Error()))
				time.Sleep(10 * time.Second) // Ждем перед повтором
				continue
			}

			// Добавляем заказы в кэш, если их там нет
			for _, ord := range orders {

				h.Cache.Delete(ord.OrderUID)
				h.Cache.Store(ord.OrderUID, ord)
			}

			h.Logger.Info("Кэш синхронизирован с базой данных")
			time.Sleep(1 * time.Minute) // Интервал между синхронизациями
		}
	}()
}

// NewAPIHandler создает новый экземпляр APIHandler
func NewAPIHandler(db *gorm.DB, cacheInstance *sync.Map, redisClient *cache.RedisClient) *APIHandler {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	return &APIHandler{
		DB:     db,
		Cache:  cacheInstance,
		Redis:  redisClient,
		Logger: logger,
	}
}

// HealthCheckHandler обрабатывает запросы на проверку состояния
func (h *APIHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	h.Logger.Info("Обработан запрос проверки состояния")
}

// GetOrderHandler обрабатывает запрос на получение заказа по ID
func (h *APIHandler) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderUID := vars["id"]

	// Проверяем наличие в кэше
	if value, ok := h.Cache.Load(orderUID); ok {
		if ord, ok := value.(order.Order); ok {
			h.Logger.Info("Данные перед отправкой из кэша", slog.Any("response", ord))
			json.NewEncoder(w).Encode(ord)
			return
		}
	}

	// Загружаем из базы данных
	var ord order.Order
	if err := h.DB.
		Preload("Items").
		Preload("Delivery").
		Preload("Payment").
		First(&ord, "order_uid = ?", orderUID).Error; err == nil {

		// Логируем данные перед кэшированием и отправкой
		h.Logger.Info("Данные перед отправкой из базы", slog.Any("response", ord))

		// Кэшируем данные и отправляем их
		h.Cache.Store(ord.OrderUID, ord)
		json.NewEncoder(w).Encode(ord)
	} else {
		h.Logger.Warn("Заказ не найден", slog.String("orderUID", orderUID))
		http.Error(w, "Заказ не найден", http.StatusNotFound)
	}
}

// DeleteOrderHandler обрабатывает запрос на удаление заказа по ID
func (h *APIHandler) DeleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderUID := vars["id"]

	// Удаляем заказ из базы данных
	if err := h.DB.Delete(&order.Order{}, "order_uid = ?", orderUID).Error; err != nil {
		h.Logger.Error("Не удалось удалить заказ из базы данных",
			slog.String("orderUID", orderUID), slog.String("error", err.Error()))
		http.Error(w, "Не удалось удалить заказ", http.StatusInternalServerError)
		return
	}

	// Удаляем заказ из кэша
	h.Cache.Delete(orderUID)
	h.Logger.Info("Заказ удален из кэша", slog.String("orderUID", orderUID))

	// Удаляем заказ из Redis
	h.Redis.Client.Del(r.Context(), orderUID)
	h.Logger.Info("Заказ удален из Redis", slog.String("orderUID", orderUID))

	// Отправляем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Заказ успешно удален"))
	h.Logger.Info("Заказ успешно удален", slog.String("orderUID", orderUID))
}
