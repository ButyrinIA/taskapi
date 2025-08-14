Task API
REST API для управления задачами на Go 1.23.10 с использованием принципов чистой архитектуры. Поддерживает асинхронное логирование действий и ошибок через канал и отдельную горутину, in-memory хранилище и graceful shutdown.
Структура

cmd/api/main.go: Точка входа, настройка сервера и graceful shutdown.
internal/adapters/handlers/: HTTP-обработчики.
internal/adapters/repositories/: Реализация in-memory репозитория.
internal/core/models/: Сущности (Task).
internal/core/usecases/: Бизнес-логика (usecases).

Функционал

GET /tasks?status=: Получить список задач, с фильтрацией по статусу (например, "todo", "in_progress", "done"), если указано.
GET /tasks/{id}: Получить задачу по ID.
POST /tasks: Создать новую задачу (JSON-тело: {"title": "...", "description": "...", "status": "..."}).
Асинхронное логирование действий и ошибок в stdout.
In-memory хранилище с потокобезопасностью (используется sync.Mutex).
Graceful shutdown при получении сигналов SIGINT/SIGTERM.

Требования

Go 1.23.10
Docker (для запуска в контейнере)
Git

Сборка и запуск

Убедитесь, что установлен Go 1.23.10.
Клонируйте репозиторий: git clone <https://github.com/ButyrinIA/taskapi>.
Перейдите в директорию проекта: cd taskapi.
Инициализируйте Go-модуль: go mod init github.com/ButyrinIA/taskapi.
Соберите приложение: go build -o taskapi cmd/api/main.go.
Запустите: ./taskapi.
Сервер доступен по адресу: http://localhost:8080.

Сборка и запуск через Docker

Убедитесь, что установлен Docker.
Соберите образ: docker build -t taskapi .
Запустите контейнер: docker run -p 8080:8080 taskapi
API доступен по адресу: http://localhost:8080.

Тестирование
Используйте инструменты, такие как curl или Postman, для проверки API:

Создание задачи:curl -X POST -H "Content-Type: application/json" -d '{"title":"Test Task","description":"Description","status":"todo"}' http://localhost:8080/tasks


Получение списка задач:curl http://localhost:8080/tasks?status=todo


Получение задачи по ID:curl http://localhost:8080/tasks/1



Для запуска модульных тестов:
go test ./... -v

Замечания
Логи (действия и ошибки) выводятся в stdout. Для продакшена можно настроить вывод в файл или внешний сервис.
Тесты покрывают репозиторий, usecases и обработчики, включая сценарии с ошибками (например, пустой заголовок задачи).

