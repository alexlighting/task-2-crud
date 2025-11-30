## Результаты ручного тестирования

#### Получение списка задач
Запрос
```bash
>curl http://localhost:8080/tasks
```
Ожидаемый статус 200  
Ответ
```bash
[{"id":1,"title":"Дефрагментация диска","done":false,"created_at":"2025-11-30T20:15:22+03:00"},{"id":2,"title":"Сборка мусора","done":false,"created_at":"2025-11-30T20:15:22+03:00"},{"id":3,"title":"Оптимизация запросов","done":false,"created_at":"2025-11-30T20:15:22+03:00"}]
```

#### Cоздание новой задачи
###### Корректный запрос
Запрос
```bash
>curl -X POST -H "Content-Type: application/json" --data "{\"title\": \"Создание индексов\", \"done\": false}" http://localhost:8080/tasks
```
Ожидаемый статус 201  
Ответ
```bash
{"id":4,"title":"Создание индексов","done":false,"created_at":"2025-11-30T20:42:07+03:00"}
```

###### Некрректный запрос
Запрос
```bash
>curl -X POST -H "Content-Type: application/json" --data "{\"done\": false}" http://localhost:8080/tasks
```
Ожидаемый статус 400  
Ответ
```bash

{"error":"Task title is empty"}

```

#### Просмотр задачи
###### Корректный запрос
Запрос
```bash
>curl http://localhost:8080/tasks/4
```
Ожидаемый статус 200  
Ответ
```bash
{"id":4,"title":"Создание индексов","done":false,"created_at":"2025-11-30T20:42:07+03:00"}
```

###### Некрректный запрос
Запрос
```bash
>curl http://localhost:8080/tasks/sef
```
Ожидаемый статус 400  
Ответ
```bash
{"error":"The input value given is not a valid integer"}
```

#### Изменение задачи
###### Корректный запрос
Запрос
```bash
>curl -X PUT -H "Content-Type: application/json" --data "{\"title\": \"Нормировка индексов\", \"done\": true}" http://localhost:8080/tasks/4
```
Ожидаемый статус 200  
Ответ
```bash
{"status":"Task succesfully updated"}
```

###### Некрректный запрос
Запрос
```bash
>curl -X PUT -H "Content-Type: application/json" --data "\"title\": \"Нормировка индексов\", \"done\": true" http://localhost:8080/tasks/4
```
Ожидаемый статус 400  
Ответ
```bash
{"error":"Incorrect JSON"}
```

#### Удаление задачи
###### Корректный запрос
Запрос
```bash
>curl -X DELETE http://localhost:8080/tasks/4
```
Ожидаемый статус 400  
Тело ответа пустое.

###### Некрректный запрос
Запрос
```bash
>curl -X DELETE http://localhost:8080/tasks/4
```
Если задача с ID 4 уже удалена или еще не была создана
Ожидаемый статус 404  
Ответ
```bash
{"error":"Wrong id"}
```

#### Обращение к неподдерживаемому методу
Запрос
```bash
>curl -X PATCH http://localhost:8080/tasks/4
```
Ожидаемый статус 405  
Ответ
```bash
{"Allow":"GET, POST, PUT, DELETE"}
```