<h3 align="center">
  <div align="center">
    <h1>Mattermost Poll Bot 🤖</h1>
  </div>
  </a>
</h3>

## Описание

**Mattermost Poll Bot** — это бот для сервера **Mattermost**, который позволяет пользователям создавать голосования, голосовать и просматривать результаты. Бот поддерживает два режима хранения данных:

- **In-Memory** (по умолчанию)
- **Tarantool** (если включено в настройках)

---

## 📋 Возможности

- Создание голосований через слеш-команды
- Голосование за предложенные варианты
- Получение результатов голосования
- Закрытие голосования
- Удаление голосования

---

## Установка

```sh
git clone https://github.com/goroutiner/matterpoll-bot.git
```

---

### 🔧 Предварительная конфигурация (для демонстрации возможностей бота)

---

1. Убедитесь, что у вас установлен **Docker** и **Docker Compose** и выполните команду для запуска **Mattermost**:

```sh
docker compose up -d mattermost
```

2. Перейдите на web-версию **Mattermost**:http://localhost:8065 и выполните следующую инструкцию [инструкция](/instructions)

3. Прописываем имя команды и полученный токен бота в переменные окружения в файле `compose.yaml` для сервиса `matterpoll-bot`:

```yaml
environment:
  MODE: "database"
  SERVER_URL: "http://mattermost:8065"
  BOT_SOCKET: ":4000"
  DB_SOCKET: "tarantool:3301"
  TEAM_NAME: "your_team_name"
  BOT_TOKEN: "your_bot_token"
```

Если необходим **memory** режим, то укажите:

```yaml
MODE: "memory"
```

---

### 🔧 Предварительная конфигурация (для личного пространства)

1. Прописываем имя команды и токен бота в переменные окружения в файле `compose.yaml` для сервиса `matterpoll-bot`:

```yaml
environment:
  MODE: "database"
  SERVER_URL: "http://your_domain"
  BOT_SOCKET: ":4000"
  DB_SOCKET: "tarantool:3301"
  TEAM_NAME: "your_team_name"
  BOT_TOKEN: "your_bot_token"
```

Если необходим **memory** режим, то укажите:

```yaml
MODE: "memory"
```

---

## 🐳 Запуск через Docker Compose

Если вы хотите запустить проект через Docker, следуйте этим шагам:

1. Убедитесь, что у вас запущен Docker
2. Перейдите в корневую директорию проекта
3. Соберите и запустите приложение с помощью команды:

- Запуск с конфигурациями **для демонстрации возможностей бота**:
```sh
docker compose up -d
```

- Запуск с конфигурациями **для личного пространства**:
```sh
docker compose up -d matterpoll-bot tarantool
```

4. Теперь, когда бот запущен, то можно использовать его функционал.

---

## 🧑🏽‍💻 Примеры использования функционала:

1. Отображение созданных команд: \
   ![Commands](https://github.com/goroutiner/matterpoll-bot/raw/main/instructions/images/commands.png)

2. Создание голосования:

```sh
/poll-create "Example" "Option1" "Option2"
```

![Created Poll](https://github.com/goroutiner/matterpoll-bot/raw/main/instructions/images/created_poll.png)

3. Получение результатов:

```sh
/poll-results "h3twm167pjgibyb5acdcjut5to"
```

![Created Poll](https://github.com/goroutiner/matterpoll-bot/raw/main/instructions/images/poll_results.png)

---

## ✅⭕ Инструкция по запуску тестов

Для запуска интеграционных тестов выполните команду:

```sh
make unit-tests
```

---

## 🛠️ Технические ресурсы

- **Язык программирования**: Go (Golang)
- **База данных**: Tarantool
- **Библиотеки**:
  - [mattermost/mattermost-server/v6](github.com/mattermost/mattermost-server/v6) для работы с API Mattermost.
  - [tarantool/go-tarantool/v2](github.com/tarantool/go-tarantool/v2) для взаимодействия с базой данных.
  - [stretchr/testify](https://github.com/stretchr/testify) для написания тестов.
