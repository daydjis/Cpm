# Telegram-бот рассылки по категориям

Бот парсит товары с сайта по расписанию и отправляет подписчикам уведомления о новых позициях. Пользователь может выбрать категории для подписки из списка.

## Требования

- Go 1.21+
- PostgreSQL

## Структура проекта

```
cmd/app/          — точка входа (main)
internal/
  bot/            — Telegram-бот (Init, обёртка handlers)
  bot/handlers/   — обработчики по папкам (start, subscribe, unsubscribe, categories, add_category, callbacks, notifications)
  config/         — загрузка конфигурации из .env
  db/             — подключение к БД
  parser/         — парсинг по slug категории (URL = SCRAPE_BASE_URL/slug)
  repository/     — слой доступа к данным
  scheduler/      — шедулер (интервал парсинга и рассылки)
  service/        — бизнес-логика
models/           — доменные модели
migrations/       — SQL-миграции схемы и сиды
```

## Настройка

1. Клонируйте репозиторий и перейдите в каталог проекта.

2. Скопируйте пример конфигурации:
   ```bash
   cp .env.example .env
   ```

3. Отредактируйте `.env`:
   - `TG_BOT_TOKEN` — токен бота от BotFather
   - `DATABASE_URL` — строка подключения к PostgreSQL (логин, пароль, хост, порт, имя БД)
   - `SCRAPE_INTERVAL` — интервал запуска парсера и рассылки (например `10m`, `1h`, `12h`)
   - `SCRAPE_BASE_URL` — базовый URL каталога (в конец подставляется slug категории: eldar, orks и т.д.)

4. Примените миграции (например, через [golang-migrate](https://github.com/golang-migrate/migrate)):
   ```bash
   migrate -path migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" up
   ```

5. Запуск:
   ```bash
   go run ./cmd/app
   ```

## Команды и UI бота

- `/start` — регистрация и главное меню (кнопки)
- `/subscribe` или кнопка «Подписаться» — выбор категории для подписки (inline-кнопки)
- `/unsubscribe` или кнопка «Мои подписки» — отписаться от категории
- `/categories` или кнопка «Все категории» — список категорий
- `/add_category` или кнопка «Добавить категорию» — добавить категорию по slug (например eldar, orks)

## Переменные окружения

| Переменная             | Описание                          | Обязательная |
|------------------------|-----------------------------------|--------------|
| `TG_BOT_TOKEN`         | Токен Telegram-бота              | Да           |
| `DATABASE_URL`         | URL подключения к PostgreSQL     | Да           |
| `SCRAPE_INTERVAL`  | Интервал парсинга (например 10m) | Нет (10m) |
| `SCRAPE_BASE_URL`  | Базовый URL каталога (…/warhammer-40000) | Нет |

## Docker

Для сборки и запуска через Docker используйте `Dockerfile` и `docker-compose.yml`. В `docker-compose` задайте переменные окружения или используйте `.env`.
