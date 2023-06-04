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

Test the application: `/ping`
```
curl 'localhost:8080/ping'
```

Send messages to a receiver: `/api/send`
```
curl --request POST \
  --data '{"Chat":"john:doe", "Text": "hi doe", "Sender": 1}' \
  http://localhost:8080/api/send
```

Receive past requests: `/api/pull`
```
curl --request GET \
  --data '{"Chat":"john:doe", "Cursor": 0, "Limit": 10, "Reverse": false}' \
  http://localhost:8080/api/pull
```
