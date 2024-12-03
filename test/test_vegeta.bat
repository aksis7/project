@echo off
REM Создаем файл с HTTP-запросом
echo GET http://wb-go-service:8082/orders/b563feb7b2b84b6test > input.txt

REM Выполняем нагрузочное тестирование через Vegeta
docker exec -i wb-vegeta vegeta attack -rate=15000 -duration=30s < input.txt | docker exec -i wb-vegeta vegeta report

REM Удаляем временный файл (опционально)
del input.txt
