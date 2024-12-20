# WB Microservices Project

## Описание проекта

Этот проект представляет собой набор микросервисов, которые включают:
- **Kafka** для обработки сообщений.
- **Zookeeper** для координации Kafka.
- **PostgreSQL** как основную базу данных.
- **Redis** для кэширования.
- **Go-приложение**, взаимодействующее с Kafka, Redis и PostgreSQL.
- Инструменты для нагрузочного тестирования: **Vegeta** и **WRK**.

Проект разворачивается с использованием Docker Compose.

---

## Сервисы

- **Zookeeper**: используется для координации Kafka.
- **Kafka**: брокер сообщений.
- **PostgreSQL**: реляционная база данных для хранения заказов.
- **Redis**: система кэширования.
- **Go-сервис**: основное приложение, предоставляющее API.
- **Vegeta** и **WRK**: инструменты для нагрузочного тестирования.

---

## Установка и запуск

1. Убедитесь, что у вас установлены:
   - **Docker**: [Инструкция по установке](https://docs.docker.com/get-docker/).
   - **Docker Compose**: [Инструкция по установке](https://docs.docker.com/compose/install/).

2. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/aksis7/project.git
   cd project
   
3. Запустите сервисы:
     ```bash
     docker-compose up --build

4.  Проверьте, что все сервисы работают:
     ```bash
     docker-compose ps

5.  Тестирование веб-интерфейса:
      ```bash
       http://localhost:8082/health
       Ответ:ок
7.  Откройте веб-интерфейс в браузере:
      ```bash
      Перейдите по адресу: http://localhost:8082.
      Введите нужный uid,например:b563feb7b2b84b6test
9.  Тестирование производительности:
     ```bash
     Запустите test/vegeta-test.bat
     Запустите  test/wrk-test.bat
11.  Остановка всех сервисов:
      ```bash
      
      docker-compose down
    
