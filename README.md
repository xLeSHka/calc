# calc
Это сервер-калькулятор. Он использует урезанный [Shunting yard algorithm](https://en.wikipedia.org/wiki/Shunting_yard_algorithm) для вычисления значений выражений. 
Выражение должно состоять из
```
digit = "0" ... "9" 
operators = "*" | "/" | "+" | "-" 
punctuation = "." | "(" | ")" 
```
Операторы и их приоритет
| Operator | Precedence |
|:--------:|:----------:|
| *        |2           |
|/         |2           |
|+         |1           |
|-         |1           |
## Запуск сервера через docker
Запустить его можно через `docker`, `make` или сбилдить самому
### Docker
Чтобы запустить через `docker` необходим собственно [docker](https://docs.docker.com/compose/install/). После установки просто введите в терминал `VS Code`, при этом `Docker Desktop` обязательно должен быть запущен
```cmd
docker-compose up -d
```
Тогда логи сервера можно посмотреть в `Docker Desktop` нажав на `calc_service`
Чтобы остановить сервер нужно ввести в терминал `VS Code` 
```cmd
docker compose down
```
## Запуск сервера без docker
Если же вы не хотите или не можете запустить сервер через `docker`, хотя я настоятельно рекомендую именно этот способ, докер очень удобен и скорее всего пригодится вам еще много раз, то можно сделать это командой `make` или просто сборкой бинарника и его запуском, но для этого нужно в `main.go`, там где мы задаем порт для сервера изменить `":%d"` на `"localhost:%d"`  
### Makefile
Теперь можно скачать [make](https://stackoverflow.com/questions/32127524/how-to-install-and-use-make-in-windows) и ввести команду `make` для сборки бинарника и его запуска
```cmd
make
``` 
### Самостоятельная сборка
Для самостоятельной сборки нужно ввести команды в терминал `VS Code`
```cmd
go build -o calc_service ./cmd/main/main.go
./calc_service
```
## Тестирование сервера
Протестирвоать сервер можно с помощью заготовленных `curl` запросов, автотестов, через `Postman` или через `swagger-ui`, если вы запускали сервер через `docker`
### Автотесты 
Для запуска тестов нужно ввести в терминал `VS Code`
```cmd
go test ./internal/server/ -v -cover
go test ./pkg/calculator/ -v -cover
```
### Curl запросы
Введите команду:
```cmd
curl -w "%{http_code}" --location 'localhost:9090/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
```
**Ожидаемый ответ:**
```json
{
   "expression":"2+2*2",
   "result":"6"
}
```
```cmd 
200
```
Введите команду:
```
curl -w "%{http_code}" --location 'localhost:9090/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": ""
}'
```
**Ожидаемый ответ:**
```json
{
   "error":"Expression is not valid"
}
```
```cmd 
422
```
Введите команду:
```
curl -w "%{http_code}" --location 'localhost:9090/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "internal"
}'
```
**Ожидаемый ответ:**
```json
{
   "error":"Internal server error"
}
```
```cmd 
500
```
### Swagger-UI
Если вы подняли этот сервер с `docker`, можно использовать [swagger-ui](http://localhost:8085/), там будет удобный интерфейс для создания своих запросов. 
### Postman
Так же можно использовать `Postman`. Если вы пользуетесь `VS Code`, то нужно просто зайти в `extention` в `VS Code`, ввести `Postman` и установить первое расширение из списка. Чтобы пользоваться `Postman` нужно в нем зарегистрироваться. После регистрации нужно зайти в свой аккаунт в расширении для `VS Code`. И все, можно создавать запросы нажатием на `NewHTTPRequest`. Потом выбрать метод, ввести `localhost:9090/api/v1/calculate` в поле `URL`. Если вы хотите проверить правильность вычислений то выбранный метод должен быть `POST` и  во вкладке `body` выбрать `raw`, а потом справа нажав на синюю стрелочку выбрать `json`. Туда нужно вставить струтуру
```json
{  
    "expression":"ваше выражение"  
}  
```






