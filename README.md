Distributed Notification Service 🚀

An enterprise-grade, highly concurrent distributed notification system built with Golang. This microservice architecture handles asynchronous email and SMS routing, featuring rate-limiting, idempotency, and fault-tolerant delivery guarantees.

🏗 Architecture

The system is fully decoupled. The API gateway validates and queues requests, while background workers consume jobs, interface with third-party providers (SMTP, SendGrid, Twilio), and handle automatic retries.

graph TD
    Client[Client Request] --> API[Go API Gateway]
    API <-->|Check Rate Limit & Idempotency| Redis[(Redis)]
    API -->|Publish (Async)| NATS[NATS JetStream Queue]
    NATS -->|Consume Jobs| Worker[Go Background Worker]
    Worker -->|Attempt 1-3| External[External Providers: SMTP / Twilio]
    External -- Failure --> Worker
    Worker -->|Attempt 4+| DLQ[Dead Letter Queue DLQ.NOTIFY]



✨ Key Features

Interface-Driven Providers: Cleanly separated integration layer supporting Gmail SMTP, SendGrid, and Twilio. Easy to extend.

At-Least-Once Delivery: Powered by NATS JetStream. If a worker crashes mid-process, the message is safely re-queued.

Dead Letter Queue (DLQ) & Retry Loops: Network failures or API rejections trigger exponential backoff. Poison-pill messages are safely routed to a DLQ after 3 failed attempts to prevent queue blocking.

Idempotency & Rate Limiting: Utilizes Redis SetNX and sliding window counters to prevent duplicate notifications and user spam.

Fully Containerized: Multi-stage Docker builds ensure the Go binaries are tiny (<20MB), and docker-compose orchestrates the entire stack effortlessly.

🛠 Tech Stack

Language: Golang (Go 1.22+)

Message Broker: NATS JetStream

Cache / State: Redis

Containerization: Docker & Docker Compose

🚀 Quick Start (Local Development)

Start the entire distributed system (API, Worker, NATS, Redis) in isolated containers:

docker compose up --build



Trigger a Notification:

curl -X POST http://localhost:8080/api/v1/notify \
-H "Content-Type: application/json" \
-d '{
    "user_id": "usr_123",
    "type": "EMAIL",
    "target": "user@example.com",
    "title": "Welcome!",
    "body": "Your account has been created successfully."
}'

