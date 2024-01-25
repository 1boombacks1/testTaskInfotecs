# Используемые технологии:
- **Go**
- **Postgres**
- **Fiber** (веб фреймворк)
- **Docker** (для запуска сервиса)
# Особенности
- Сервис был написан с Clean Architecture, что позволяет легко расширять функционал сервиса и тестировать его.
- Также был реализован Graceful Shutdown для корректного завершения работы сервиса.
- Присутствют тесты слоя репозитория
# Usage
#### Запуск сервиса
Команда `make compose-up` - запускает сервис
#### Остановка сервиса
Команда `make compose-down` - останавливает сервис