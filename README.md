
# Web-calculator | [English](README.eng.md) | [Русский](README.md)

**Web-calculator** – это распределённая система, реализующая обработку арифметических выражений. В основе системы лежит взаимодействие двух основных компонентов:

- **Оркестратор**
  Центральный сервер, который:
  - Предоставляет вычислительные задачи посредством ```GET``` запроса к эндпоинту ```/internal/task```
  - Принимает результаты вычислений посредством ```POST``` запроса к тому же эндпоинту ```/internal/task```

- **Агент**  
  Компонент, который:
  - Периодически запрашивает у оркестратора задачу "дай задачку поработать"
  - Выполняет вычисления
  - Отправляет вычисленный результат обратно на оркестратор

### Архитектура
<img src="./img/image.png" alt="Схема проекта" width="1000">

## Функционал

- Поддерживаются операции сложения, вычитания, умножения и деления, а также выражения в скобках
- Выражение может вводиться как с пробелами между числом и операндом, так и без
- Калькулятор принимает на вход положительные целые числа


## Зависимости

- Go версии ```1.23``` или новее
- Дополнительные библиотеки (указаны в ```go.mod```)

## ⚠️ Предупреждение
**Уважаемые проверяющие ЯЛ**  
Комментарии в коде от и до написаны моими руками и удовлетворяют требованиям ```go-lint```, который я использую в своем проекте для соблюдения код стайла

## Установка

1. Клонирование репозитория

```bash
git clone https://github.com/bulbosaur/web-calculator-golang
```

2. Запуск сервера из репозитория проекта

Необходимо находиться в корневой директории проекта (web-calculator-golang)

Для запуска каждого из двух сервисов Вам потребуется 2 отдельных терминала. Удобно открыть сплит терминалов можно с помощью горячих клавиш: ```Ctrl + Shift + 5``` (Windows/Linux) или ```Cmd + Shift + 5``` (macOS)

В первом необходимо ввеси команду:

```bash
go run ./cmd/orchestrator/main.go
```

А во втором:

```bash
go run ./cmd/agent/main.go
```


## Переменные окружения

| Переменная                    | Описание                                            | Значение по умолчанию |
|-------------------------------|-----------------------------------------------------|-----------------------|
| ```PORT```                    | Порт для запуска сервера                            | 8080                  |
| ```HOST```                    | Хост для запуска сервера                            | localhost             |
| ```TIME_ADDITION_MS```        | Время выполнения операции сложения в миллисекундах  | 100                   |
|```TIME_SUBTRACTION_MS```      | Время выполнения операции вычитания в миллисекундах | 100                   |
| ```TIME_MULTIPLICATIONS_MS``` | Время выполнения операции умножения в миллисекундах | 100                   |
| ```TIME_DIVISIONS_MS```       | Время выполнения операции деления в миллисекундах   | 100                   |
| ```DATABASE_PATH```           | Путь к базе данных                                  |                       |


Чтобы изменить значения переменных окружения, необходимо создать файл ```config.yaml``` (или отредактировать существующий файл ```example_config.yaml```)

### Как должен выглядеть config файл

```bash
# web-calculator-golang\config\config.yaml
server:
  host: localhost
  port: 8080

time:
  TIME_ADDITION_MS: 100
  TIME_SUBTRACTION_MS: 100
  TIME_MULTIPLICATIONS_MS: 100
  TIME_DIVISIONS_MS: 100

worker:
  COMPUTING_POWER: 5
  
database:
  DATABASE_PATH: ./db/calc.db
```

## Публичный API

Базовый URL по умолчанию: ```http://localhost:8080```

| API endpoint             | Метод      | Тело запроса                  | Ответ сервера                                                        | Код ответа |
|-------------------------|-------------|-------------------------------|----------------------------------------------------------------------|------------|
| ```/api/v1/calculate``` | ```POST```  | ```{"expression": "2 * 2"}``` | ```{"id":1}```                                                       | 200        |
| ```/api/v1/calculate``` | ```POST```  | ```"expression": "2 * 2"```   | ```{"error":"Bad request","error_message":"invalid request body"}``` | 400        |
| ```/api/v1/calculate``` | ```GET```   | ```{"expression": "2 * 2"}``` | ```Method Not Allowed```                                             | 405        |
| ```/coffee```           |             |                               | ```I'm a teapot```                                                   | 418        |
| ```/api/v1/tea```       |             |                               |```404 page not found```                                              | 404        |

### Коды ответов

- 200 - Успешный запрос
- 400 - Некорректный запрос
- 404 - Ресурс не найден
- 405 - Метод не поддерживается 
- 422 - Некорректное выражение (например, буква английского алфавита вместо цифры)
- 500 - Внутренняя ошибка сервера

### Примеры работы

Для отправки POST запросов удобнее всего использовать программу [Postman](https://www.postman.com/downloads/).

## База данных

<img src="https://i.imgur.com/CDJrb9i.png" alt="Схема БД" width="1500">

## Тестирование

```bash
go test -v  ./...
```