@echo off
REM Запуск команды wrk с параметрами нагрузки

echo Запуск wrk с параметрами нагрузки...
docker exec -i wb-wrk wrk -t12 -c400 -d30s http://wb-go-service:8082/orders/test_order_1

REM Проверка завершения работы
IF %ERRORLEVEL% EQU 0 (
    echo Команда выполнена успешно.
) ELSE (
    echo Произошла ошибка при выполнении команды.
)

pause
