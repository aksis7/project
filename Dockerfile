# Dockerfile

# Используем golang как базовый образ для сборки
FROM golang:1.23.3 AS builder

# Создаем директорию для приложения
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь исходный код в контейнер
COPY . .
COPY static /app/static
# Переходим в директорию с main.go и собираем бинарник
WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .

# Используем минимальный образ для запуска
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/cmd/main .
RUN chmod +x main
# Порт, который будет использоваться приложением
EXPOSE 8082

# Запускаем собранное приложение
CMD ["./main"]
