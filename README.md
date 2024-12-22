
# web-calculator

Web-calculator представляет из себя веб-сервис, при помощи которого пользователь может отправить арифметическое выражение по HTTP и получить в ответ его результат.



## Функционал

- Поддерживаются операции сложения, вычитания, умножения и деления
- Выражение может вводиться как с пробелами, так и без
- Калькулятор принимает на вход положительные целые числа


## Требования

- Go версии ```1.22``` или новее


## Установка

1. Клонирование репозитория

```bash
git clone https://github.com/bulbosaur/web-calculator-golang
```

2. Запуск сервера
   
```bash
cd web-calculator-golang
go run cmd/main.go
```


## Переменные окружения

| Переменная | Описание | Значение по умолчанию |
|------------|----------|-----------------------|
| ```PORT``` | Порт для запуска сервера | ```8080``` |
| ```HOST``` | Хост для запуска сервера | ```localhost``` |



## API

Базовый URL по умолчанию: ```http://localhost:8080```

| API endpoint | Метод | Тело запроса | Ответ сервера | Код ответа |
|--------------|-------|--------------|---------------|------------|
| ```/api/v1/calculate``` | ```POST``` | ```{"expression": "2 * 2"}``` | ```{"result":"4"}``` | 200 |
| ```/api/v1/calculate``` | ```POST``` | ```"expression": "2 * 2"``` | ```{"error":"Bad request","error_message":"invalid request body"}``` | 400 |
| ```/api/v1/calculate``` | ```GET``` | ```{"expression": "2 * 2"}``` | ```Method Not Allowed``` | 405 |
| ```/coffee``` | | | | ```I'm a teapot``` | 418 |
| ```/api/v1/tea``` | | | | ```404 page not found``` | 404 |

### Коды ответов

- 200 - Успешный запрос
- 400 - Некорректный запрос
- 404 - Ресурс не найден
- 405 - Метод не поддерживается 
- 422 - Некорректное выражение (например, буква английского алфавита вместо цифры)
- 500 - Внутренняя ошибка сервера

### Примеры работы

Для отправки POST запросов удобнее всего использовать программу [Postman](https://www.postman.com/downloads/).

1. StatusOK 200
```bash
curl -x POST --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "42 + 5 * 2"
}'

# {"result":52}
```

```bash
curl -x POST --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "6-8"
}'

# {"result":-2}
```

```bash
curl -x POST --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "123(3/2)"
}'

# {"result":184.5}
```

2. Bad Request 400

```bash
curl -x POST --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2 * 2
}'

# {"error":"Bad request","error_message":"invalid request body"}
```

3. Unprocessable Entity 422
```bash
curl -x POST --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "cat + 100500"
}'

# {"error":"Expression is not valid","error_message":"invalid characters in expression"}
```

```bash
curl -x POST --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "()"
}'

# {"error":"Expression is not valid","error_message":"the brackets are empty"}
```

```bash
curl -x POST --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "1/0"
}'

# {"error":"Expression is not valid","error_message":"division by zero is not allowed"}
```