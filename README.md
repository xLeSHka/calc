# calc
Это сервер-калькулятор. Он использует [Shunting yard algorithm](https://en.wikipedia.org/wiki/Shunting_yard_algorithm) для вычисления значений выражений. 
Выражение должно состоять из
```
digit = "0" ... "9" 
operators = "*" | "/" | "+" | "-" | "^"
punctuation = "." | "," | "(" | ")" 
function = log(a,x) | sqrt(a)
```
Операторы и их приоритет
| Operator | Precedence |
|:--------:|:----------:|
| ^        |5           |
|-(unary)  |4           |
| *        |3           |
|/         |3           |
|+         |2           |
|-         |2           |
  
**Unary минус** нужно выделять в скобки, если он не стоит после скобки. `log(a,x)` представляет из себя `log10(a)`/`log10(x)` то есть `loga(x)`. `sqrt(a)` это и есть корень квадратный из а, думаю тут объяснять не нужно. На этом все, мне было лень добавлять больше функций :0
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
**Запросы вводить нужно не в `cmd` или `power shell`, а в `Git Bash`, если вы скачивали `Git` себе на компьютер, то он у вас должен быть. Так как в `cmd` и `power shell` нужно экранировать например `^`. В итоге читать и писать выражения становится в разы труднее**  
Правильный запрос  
Введите команду:
```cmd
curl -X 'POST' -w "%{http_code}"\
  'http://localhost:9090/api/v1/calculate' \
 -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
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
Правильный запрос  
Введите команду:
```cmd
curl -X 'POST' -w "%{http_code}"\
  'http://localhost:9090/api/v1/calculate' \
 -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "expression": "log(18,18)^(-9)/3.14*(-12-3)*3/10+2*sqrt(4)"
}'
```
**Ожидаемый ответ:**
```json
{
   "expression":"log(18,18)^(-9)/3.14*(-12-3)*3/10+2*sqrt(4)",
   "result":"2.56689"
}
```
```cmd 
200
```
Правильный запрос  
Введите команду:
```cmd
curl -X 'POST' -w "%{http_code}"\
  'http://localhost:9090/api/v1/calculate' \
 -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "expression": "log(sqrt(4),sqrt(64))"
}'
```
**Ожидаемый ответ:**
```json
{
   "expression":"log(sqrt(4),sqrt(64))",
   "result":"3"
}
```
```cmd 
200
```
Правильный запрос  
Введите команду:
```cmd
curl -X 'POST' -w "%{http_code}"\
  'http://localhost:9090/api/v1/calculate' \
 -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "expression": "log(sqrt(4),sqrt(64))^sqrt(81)*(-1)"
}'
```
**Ожидаемый ответ:**
```json
{
   "expression":"log(sqrt(4),sqrt(64))^sqrt(81)*(-1)",
   "result":"-19683"
}
```
```cmd 
200
```
Унарный минус не выделен скобками  
Введите команду:
```cmd
curl -X 'POST' -w "%{http_code}"\
  'http://localhost:9090/api/v1/calculate' \
 -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "expression": "log(sqrt(4),sqrt(64))^sqrt(81)*-1"
}'
```
**Ожидаемый ответ:**
```json
{
   "messsage":"Expression is not valid"
}
```
```cmd 
422
```
Деление на ноль  
Введите команду:
```cmd
curl -X 'POST' -w "%{http_code}"\
  'http://localhost:9090/api/v1/calculate' \
 -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "expression": "log(sqrt(4),sqrt(64))/0"
}'
```
**Ожидаемый ответ:**
```json
{
   "messsage":"Expression is not valid"
}
```
```cmd 
422
```
Пропущена закрывающая скобка   
Введите команду:
```cmd
curl -X 'POST' -w "%{http_code}"\
  'http://localhost:9090/api/v1/calculate' \
 -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "expression": "log(sqrt(4),sqrt(64)"
}'
```
**Ожидаемый ответ:**
```json
{
   "messsage":"Expression is not valid"
}
```
```cmd 
422
```
Неправильный синтаксис выражения  
Введите команду:
```cmd
curl -X 'POST' -w "%{http_code}"\
  'http://localhost:9090/api/v1/calculate' \
 -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "expression": "2*2*"
}'
```
**Ожидаемый ответ:**
```json
{
   "messsage":"Expression is not valid"
}
```
```cmd 
422
```
Запрос с неправильным Content-Type`ом  
Введите команду:
```
curl -X 'POST' -w "%{http_code}"\
  'http://localhost:9090/api/v1/calculate' \
 -H 'accept: text/plain' \
  -H 'Content-Type: text/plain' \
  -d '{
  "expression": "2+2+2"
}'
```
**Ожидаемый ответ:**
```json
{
   "message":"Expression is not valid"
}
```
```cmd 
422
```
Запрос с недопустимым методом  
Введите команду:
```cmd
curl -X 'GET' -w "%{http_code}"\
  'http://localhost:9090/api/v1/calculate' \
 -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "expression": "2+2*2"
}'
```
**Ожидаемый ответ:**
```json
{
   "message":"Method Not Allowed"
}
```
```cmd 
405
```
Заглушка для 500 ошибки  
Введите команду:
```
curl -X 'POST' -w "%{http_code}"\
  'http://localhost:9090/api/v1/calculate' \
 -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "expression": "internal"
}'
```
**Ожидаемый ответ:**
```json
{
   "message":"Internal server error"
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
