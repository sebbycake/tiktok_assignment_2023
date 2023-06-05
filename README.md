# TikTok Tech Immersion Programme 2023 Server Assignment

### Features

- Send messages to receiver
- Receive past messages

### Built With

* [Go](https://go.dev/) - The back-end language used
* [Kitex](https://github.com/cloudwego/kitex) - Golang RPC framework
* [Docker](https://www.docker.com/) - Container tool

## Getting Started

### Pre-requisites

* [Go](https://go.dev/)
* [Docker Desktop](https://www.docker.com/products/docker-desktop/)

Please download the pre-requisites and open Docker Desktop application before proceeding.

### Installation

1. Clone the project to your local machine:

```
git clone https://github.com/sebbycake/tiktok_assignment_2023.git
cd tiktok_assignment_2023
```

2. Start docker containers:
```
docker-compose -f docker-compose.yml up -d --build
```

To shutdown:
```
docker compose down
```

### Usage


Test connection with server application: `/ping`
```
curl -X GET 'localhost:8080/ping'
```

Send messages to a receiver: `/api/send`
```
curl -X POST http://localhost:8080/api/send \
    -d '{"chat":"john:doe", "text":"hi doe", "sender": "1"}' \
    -H "Content-Type: application/json"
```

Receive past requests: `/api/pull`
```
curl -X GET http://localhost:8080/api/pull \
    -d '{"chat":"john:doe", "cursor": 0, "limit": 10, "reverse": false}' \
    -H "Content-Type: application/json"
```

For more information, please refer to the API documention below.

# API Documentation

This document provides an overview of the endpoints and request bodies for the API.

## `POST /api/send`

Use this endpoint to send messages to a receiver.

### Request Body

The request body should be in JSON format and contain the following properties:

| Property    | Type   | Required | Description                  | Example |
|--------------|--------|----------|------------------------------| ----- |
| `chat`    | string | Yes      | The chat identifier between two people. Format is `<sender>:<receiver>`. | peter:tom  | 
| `text`  | string | Yes      | The message content. | Hi Tom! How's your day? |
| `sender`  | string | Yes      | The sender identifier.| 1 |

#### Example

```json
{
    "chat": "peter:tom",
    "text": "Hi Tom! How's your day?",
    "sender": "1"
}
```
### Response
The response returns the id of the row in the database.
```
1
```

---

## `GET /api/pull`

Use this endpoint to retrieve messages from a specific chat.

### Request Body

The request body should be in JSON format and contain the following properties:

| Property    | Type   | Required | Description                        | Example |
|--------------|--------|----------|-----------------------------------| ----- |
| `chat`    | string | Yes      | The chat identifier between two people. Format is `<sender>:<receiver>`. | peter:tom  | 
| `cursor`  | number | Yes      | The starting position of message's send_time in microseconds, inclusively. Default is 0. | 168324550717297 |
| `limit`  | number | Yes      | The maximum number of messages returned per request.Default is 10.| 15 |
| `reverse`  | boolean | Yes      | If true, messages are sorted in descending order by the send time. | false | 

#### Example

```json
{
    "chat": "peter:tom",
    "cursor": 0,
    "limit": 10,
    "reverse": false
}  
```

### Response Body

The response body contain the following properties:

| Property    | Type   | Description                                                  | Example |
|--------------|--------|-------------------------------------------------------------| ----- |
| `messages`    | string | The list of messages of a specific chat identifier. | peter:tom  | 
| `has_more`  | boolean | Returns true if there are more messages to be retrieved beyond the limit given in the request. | true |
| `next_cursor`  | number | The starting position (in microseconds) of the next page of messages if `has_more` is true. | 1685965073721991 |

#### Example

```json
{
    "messages": [
        {
            "chat": "peter:tom",
            "text": "test msg 1",
            "sender": "1",
            "send_time": 1685965071720997
        },
        {
            "chat": "peter:tom",
            "text": "test msg 2",
            "sender": "1",
            "send_time": 1685965076130128
        },
        {
            "chat": "peter:tom",
            "text": "test msg 3",
            "sender": "1",
            "send_time": 1685965078830078
        },
    ]
}
```
Note that `has_more` and `next_cursor` are not returned because there are only 3 messages to be retrieved, which is less than the default limit of 10.
