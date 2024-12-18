
# web-calculator

Web-calculator представляет из себя веб-сервис, позволяющий пользователю отправлять cURL запросы с арифметическими выражениями и получать в ответ его результат.



## Функционал

- Поддерживаются операции сложения, вычитания, умножения и деления
- Выражение может вводиться как с пробелами, так и без
- Поддерживаются выражения, состоящие из натуральных чисел


## Требования

- Go версии ```1.20``` или новее


## Установка

1. Клонирование репозитория

```bash
git clone https://github.com/username/calculator
cd calculator
```
2. Установка зависимостей
```bach
go mod init calculator
go mod tidy
```

3. Запуск сервера
```bach
go run cmd/main.go
```


## Переменные окружения

| Переменная | Описание | Значение по умолчанию |
|------------|----------|----------------------|
| PORT | Порт для запуска сервера | 8080 |

