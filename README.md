# Go-Kafka

## Introduction

This project demonstrates how to use Apache Kafka with Go, leveraging the Sarama library by IBM. Sarama is a stable Kafka client for GoLang, making it suitable for production use.


## Setup

### Prerequisites

- Docker
- Docker Compose
- Go (version 1.16+)

### Installation

1. **Clone the repository**:
   ```sh
   git clone https://github.com/vibzz2003/GO-KAFKA.git
   cd GO-KAFKA

2. **Start Kafka and Zookeeper**:
    Run the following command in the root directory to start Apache Kafka and Zookeeper:
    ```sh
    docker-compose up -d
    
3. Install dependencies:
   ```sh
   cd producer
   go mod tidy
   cd ../worker
   go mod tidy

4. Start the producer (REST API):
   ```sh
   cd producer
   go run producer.go

5. Start the consumer (worker):
   ```sh
   cd worker
   go run worker.go


Use the following curl commands to send test messages:
```sh
curl --location --request POST 'http://0.0.0.0:3000/api/v1/comments' \
--header 'Content-Type: application/json' \
--data-raw '{ "text": "message 1" }'

curl --location --request POST 'http://0.0.0.0:3000/api/v1/comments' \
--header 'Content-Type: application/json' \
--data-raw '{ "text": "message 2" }'



