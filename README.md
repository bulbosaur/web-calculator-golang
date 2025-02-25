
# Web-calculator | [English](README.eng.md) | [Русский](README.md)

Web-calculator представляет из себя веб-сервис, при помощи которого пользователь может отправить арифметическое выражение по HTTP и получить в ответ его результат.



## Функционал

- Поддерживаются операции сложения, вычитания, умножения и деления, а также выражения в скобках
- Выражение может вводиться как с пробелами между числом и операндом, так и без
- Калькулятор принимает на вход положительные целые числа


## Требования

- Go версии ```1.22``` или новее

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

## API

Базовый URL по умолчанию: ```http://localhost:8080```

| API endpoint             | Метод      | Тело запроса                  | Ответ сервера                                                        | Код ответа |
|-------------------------|-------------|-------------------------------|----------------------------------------------------------------------|------------|
| ```/api/v1/calculate``` | ```POST```  | ```{"expression": "2 * 2"}``` | ```{"id":1}```                                                       | 200        |
| ```/api/v1/calculate``` | ```POST```  | ```"expression": "2 * 2"```   | ```{"error":"Bad request","error_message":"invalid request body"}``` | 400        |
| ```/api/v1/calculate``` | ```GET```   | ```{"expression": "2 * 2"}``` | ```Method Not Allowed```                                             | 405        |
|```/internal/task```     |```POST```   |                               |```{"task": {"ID": 3, "ExpressionID": 2, "Arg1": 2, "Arg2": 2, "PrevTaskID1": 0, "PrevTaskID2": 0, "Operation": "*", "Status": "awaiting processing" "Result": 0,}}``` | 200 |
|```/internal/task/1```   |```GET```    |                               |```{"task": {"ID": 1, "ExpressionID": 1, "Arg1": 2, "Arg2": 3, "PrevTaskID1": 0, "PrevTaskID2": 0, "Operation": "*", "Status": "successfully done", "Result": 6,}}```| 200 |
| ```/coffee```           |             |                               | ```I'm a teapot```                                                    | 418       |
| ```/api/v1/tea``` | | | ```404 page not found``` | 404 |

### Коды ответов

- 200 - Успешный запрос
- 400 - Некорректный запрос
- 404 - Ресурс не найден
- 405 - Метод не поддерживается 
- 422 - Некорректное выражение (например, буква английского алфавита вместо цифры)
- 500 - Внутренняя ошибка сервера

### Примеры работы

Для отправки POST запросов удобнее всего использовать программу [Postman](https://www.postman.com/downloads/).
