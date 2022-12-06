
# Task Queue Asynq

Inspired by Sidekiq from Ruby, I try to make Simple Task App for Queue Processing with the help of Asynq library and Asynqmon to Serve Monitoring Dashboard. Written with Go language.

## Requirements

- Redis
- Go >= 1.18
- Kraken.io API Key and Secret

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`REDIS_URL`

`KRAKEN_API_KEY`

`KRAKEN_API_SECRET`

`PORT` <= Optional

## Precondition

_Before run this app, you must create an account at https://kraken.io/ then copy your API Key and Secret._

## Run Locally

Clone the project

```bash
  git clone https://github.com/pengdst/task-queue-asynq.git
```

Or Using SSH

```bash
  git clone git@github.com:pengdst/task-queue-asynq.git
```

Go to the project directory

```bash
  cd task-queue-asynq
```

Create environment file

```bash
  cp .env.example .env 
```

Install dependencies

```bash
  go mod download
  go mod tidy
```

Start the server

```bash
  go run task-queue-asynq/cmd/client
  go run task-queue-asynq/cmd/server
```

### Using docker-compose

Go to the project directory

```bash
  cd task-queue-asynq
```

Create docker-compose.yml file and fill environment 

```bash
  cp docker-compose.yml.example docker-compose.yml 
```

Start the server

```bash
  docker-compose up
```

## API Reference

#### Show Dashboard

```http
  GET /dashboard
```

#### Email Delivery
Takes user_id and template_id and return message queue task status.

```http
  POST /email/delivery
```

| Parameter      | Type     | Description                       |
| :------------- | :------- | :-------------------------------- |
| `user_id`      | `int`    | **Required**.                     |
| `template_id`  | `string` | **Required**.                     |

##### Response


```http
{
  "message": "enqueued tasks: id=9130232c-4ebc-4fe9-b8ef-f0ead0c0ec01 queue=low"
}
```

#### Image Resize
Takes image_url and return message queue task status.

```http
  POST /image/resize
```

| Parameter      | Type     | Description                       |
| :------------- | :------- | :-------------------------------- |
| `image_url`    | `string` | **Required**.                     |

##### Response


```http
{
  "message": "enqueued tasks: id=9130232c-4ebc-4fe9-b8ef-f0ead0c0ec01 queue=low"
}
```


## Authors

- Github [@pengdst](https://www.github.com/pengdst)
- Gitlab [@pengdst](https://www.gitlab.com/pengdst)
- Bitbucket [@pengdst](https://www.bitbucket.com/pengdst)
- LinkedIn [@pengdst](https://www.linkedin.com/in/pengdst/)


## Acknowledgements

 - [📬 Asynq: simple, reliable & efficient distributed task queue for your next Go project](https://dev.to/koddr/asynq-simple-reliable-efficient-distributed-task-queue-for-your-next-go-project-4jhg)
 - [Simple, reliable & efficient distributed task queue in Go](https://github.com/hibiken/asynq)
 - [Web UI for monitoring & administering Asynq task queue](https://github.com/hibiken/asynqmon)

## Tech Stack

**API Client:** [Resty](https://github.com/go-resty/resty)

**Server:** [Gin](https://github.com/gin-gonic/gin)

**Image Resize:** [Kraken](https://kraken.io)

## Contributing

Contributions are always welcome!