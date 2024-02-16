# Распределенный вычислитель арифметических выражений

## Деплой

Настроить переменные среды можно изменив файл конфигурации:
- [docker-compose.yml](https://github.com/f1rezy/distributed-calculator/blob/main/docker-compose.yml)

Создайте образы и поднимите контейнеры:

```sh
$ docker compose up --build
```

## API

### Получение выражения по его идентификатору

#### Request

```http
GET http://localhost:8080/expression?id=1
```

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `id` | `int` | **Required**. идентификатор выражения |

#### Response

```javascript
{
    "ID": 1,
    "Expression": "2+2*2",
    "Result": "6",
    "Status": "ok",
    "CreatedAt": "2024-02-16T17:48:15.446296Z",
    "EvaluatedAt": "2024-02-16T17:48:15.450403Z"
}
```

### Добавление вычисления арифметического выражения

#### Request

```http
POST http://localhost:8080/expression
```

```javascript
{
    "expression": "2+2*2"
}
```

#### Response

```
Приняли к обработке
id: 1
```

### Получение списка выражений со статусами

#### Request

```http
GET http://localhost:8080/expressions
```

#### Response

```javascript
[
    {
        "ID": 1,
        "Expression": "2+2*2",
        "Result": "6",
        "Status": "ok",
        "CreatedAt": "2024-02-16T17:48:15.446296Z",
        "EvaluatedAt": "2024-02-16T17:48:15.450403Z"
    },
    {
        "ID": 2,
        "Expression": "2+2*2",
        "Result": "6",
        "Status": "ok",
        "CreatedAt": "2024-02-16T17:52:01.050951Z",
        "EvaluatedAt": "2024-02-16T17:52:01.052477Z"
    }
]
```

### Получение списка доступных операций со временем их выполения

#### Request

```http
GET http://localhost:8080/operations
```

#### Response

```javascript
[
    {
        "operator": "+",
        "execution_time": 10
    },
    {
        "operator": "-",
        "execution_time": 10
    },
    {
        "operator": "*",
        "execution_time": 10
    },
    {
        "operator": "/",
        "execution_time": 10
    }
]
```

### Настройка времени выполения операций

#### Request

```http
POST http://localhost:8080/operations
```

```javascript
{
    "add_time": 10,
    "sub_time": 10,
    "mul_time": 10,
    "div_time": 10
}
```