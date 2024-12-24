# calc_service

## Описание проекта
Данный проект реализует **HTTP-сервер**, предназначенный для вычисления математических выражений. 
Сервер обрабатывает POST-запросы, содержащие выражение в формате JSON, затем вычисляет результат и возвращает его.

### Основные функции:
- Обработка арифметических выражений.
- Обработка ошибок:
  - 422 (Unprocessable Entity)
  - 500 (Internal Server Error).

---

## Примеры использования
### Успешное выполнение:
```bash
curl --location 'http://127.0.0.1:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
```
Ответ:
```
{
  result: "6.000000"
}
```

### Ошибка 422 (invalid expression):
```bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*-"
}'
```
Ответ:
```
{
  error: "invalid expression"
}
```
### Ошибка 500 (Internal Server Error):
```bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "Hello world!"
}'
```
Ответ:
```
{
error: "It's not a bug. It's a feature"
}
```

---

## Инструкция по запуску
1. Убедитесь, что у вас установлен Go:
```
go version
```
2. Если ваш проект размещен на GitHub, вы можете скачать его, используя команду git clone. Допустим, ваш репозиторий доступен по адресу https://github.com/user/repo. Вот как это сделать:

1) Откройте терминал (или командную строку).

2) Выполните команду для клонирования репозитория:

```bash
git clone https://github.com/user/repo.git
```

   Что произойдет: Команда скачает весь содержимый код из репозитория в новую папку на вашем компьютере. Эта папка будет называться так же, как репозиторий (например, repo).

3) Дождитесь завершения операции. После этого проект появится в новом каталоге, например, repo.

   
3. Перейдите в основную папку проекта.
   
4. Запустите HTTP-сервер при помощи команды:
```bash
go run .\cmd\main.go
```
