package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"wb/internal/api"
	"wb/internal/cache"
	"wb/internal/db"
	"wb/internal/kafka"
	"wb/internal/order"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

func restoreCache(apiHandler *api.APIHandler) {
	var orders []order.Order
	err := apiHandler.DB.Preload("Items").Find(&orders).Error
	if err != nil {
		log.Error("Failed to load orders from database", slog.String("error", err.Error()))
		return
	}

	for _, ord := range orders {
		apiHandler.Cache.Store(ord.OrderUID, ord)
	}

	log.Info("Cache restored from database", slog.Int("count", len(orders)))
}

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	dsn := "host=wb-postgres-1 user=user password=password dbname=orders_db port=5432 sslmode=disable"
	database := db.InitDB(dsn)

	redisClient := cache.NewRedisClient("wb-redis-1:6379")

	var cacheInstance sync.Map
	apiHandler := api.NewAPIHandler(database, &cacheInstance, redisClient)

	// Восстанавливаем кэш из базы данных
	restoreCache(apiHandler)

	// Включаем периодическую синхронизацию кэша
	apiHandler.SyncCacheWithDB()

	// Запуск Kafka Consumer
	go kafka.ConsumeKafkaMessages(ctx, database, redisClient)

	r := mux.NewRouter()

	// Маршруты API
	r.HandleFunc("/health", apiHandler.HealthCheckHandler).Methods("GET")
	r.HandleFunc("/orders/{id}", apiHandler.GetOrderHandler).Methods("GET")
	r.HandleFunc("/orders/{id}", apiHandler.DeleteOrderHandler).Methods("DELETE")

	staticDir := "/app/static"
	log.Info("Serving static files", slog.String("directory", staticDir))

	// Маршрут для корня возвращает index.html
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, staticDir+"/index.html")
	})

	// Маршруты для статических файлов
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(staticDir))))

	// Добавляем поддержку CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},             // Укажите домены, которые можно разрешить (например, "http://localhost:3000")
		AllowedMethods:   []string{"GET", "DELETE"}, // Методы, которые разрешены
		AllowedHeaders:   []string{"Content-Type"},  // Разрешённые заголовки
		AllowCredentials: true,                      // Если нужно разрешить отправку авторизации через cookie
	}).Handler(r)

	go func() {
		log.Info("Starting server", slog.Int("port", 8082))
		if err := http.ListenAndServe(":8082", corsHandler); err != nil {
			log.Error("Server failed", slog.String("error", err.Error()))
		}
	}()

	<-stop

	log.Info("Shutting down gracefully...")
	cancel()

	sqlDB, err := database.DB()
	if err != nil {
		log.Error("Error getting database connection", slog.String("error", err.Error()))
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Error("Error closing database connection", slog.String("error", err.Error()))
	}
}
