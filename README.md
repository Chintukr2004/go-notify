# 🚀 Distributed Notification System (Microservice Architecture)

An enterprise-grade, highly concurrent distributed notification platform built in Golang. This decoupled architecture routes asynchronous multi-channel notifications (Email, SMS) while guaranteeing at-least-once delivery, idempotency, rate-limiting, and automated Dead Letter Queue (DLQ) processing for third-party outages.

---

## 🏗 System Architecture

The system decouples synchronous client requests from slow third-party delivery APIs using **NATS JetStream**. A lightweight API Gateway intercepts traffic, enforces rate limits and deduplication via **Redis**, and streams background jobs to resilient worker nodes.

```mermaid
graph TD
    Client[Client / REST Request] --> API[Go API Gateway :8080]
    API <-->|Rate Limiting & Idempotency Check| Redis[(Redis Cache :6379)]
    API -->|Publish Job Async| NATS[NATS JetStream Broker :4222]
    
    subgraph Background Processing
        NATS -->|Consume Event Stream| Worker[Go Background Worker]
        Worker -->|Attempts 1 to 3| Switchboard{Provider Interface}
        Switchboard -->|SMTP / SendGrid| Email[Email Service]
        Switchboard -->|Twilio REST| SMS[SMS Service]
        Switchboard -->|Mock / Fallback| Mock[Mock Logging]
    end

    Switchboard -- Delivery Failure / Network Error --> Worker
    Worker -->|Attempt 4+ Trip Breaker| DLQ[Dead Letter Queue DLQ.NOTIFY]
    ✨ Key Engineering FeaturesDecoupled Provider Interface: Adheres strictly to dependency injection and interface-driven design. Swapping between local SMTP, SendGrid, or Twilio requires zero code changes—only environment configuration updates.Guaranteed At-Least-Once Delivery: Powered by NATS JetStream persistent storage. If a worker process dies mid-execution, unacknowledged messages (msg.Nak()) are automatically re-delivered to surviving workers.Fault-Tolerant Dead Letter Queue (DLQ): To prevent poison-pill payloads or third-party outages (e.g., SendGrid/Twilio API downtime) from blocking the main stream, failed jobs undergo exponential backoff retries. On the 4th failure, messages are cleanly intercepted, acknowledged on the main queue, and published to DLQ.NOTIFY.Idempotency & Rate Limiting: Redis SetNX locks and sliding-window expiration counters prevent duplicate notification dispatching and protect downstream APIs from client traffic spikes.Lightweight Containerization: Compiled via multi-stage Alpine Docker builds, keeping final Go binaries under 20 MB with sub-second container startup times.🛠 Tech Stack & DecisionsTechnologyRoleWhy It Was ChosenGolang (Go 1.22+)Core ServicesHigh concurrency (goroutines), minimal memory footprint, and compile-time type safety.NATS JetStreamMessage BrokerLightweight binary transport, native persistent streaming, and ultra-low latency compared to heavy JVM brokers.Redis (Alpine)State / CachingIn-memory atomic locks (SetNX) make deduplication and rate-limiting instantaneous.Docker & ComposeOrchestrationMulti-stage container builds allow recruiters and engineers to launch the full 4-node cluster in one command.🚀 Quick Start (Recruiter / Evaluation Guide)You can spin up the entire multi-container infrastructure (API Gateway, Background Worker, NATS Broker, and Redis Cache) locally in seconds.1. PrerequisitesDocker & Docker Compose installed on your system.2. Environment SetupCreate a .env file in the root directory (or rename .env.example). To run entirely locally without third-party API keys, use the mock or smtp providers:Code snippet# System Configuration
REDIS_ADDR=redis:6379
NATS_URL=nats://nats:4222

# Provider Selection (Options: mock, smtp, sendgrid)
EMAIL_PROVIDER=mock

# Optional: Real Gmail SMTP Testing
GMAIL_ADDRESS=your.email@gmail.com
GMAIL_APP_PASSWORD=your16characterapppassword
3. Launch the StackRun the following command from the repository root:Bashdocker compose up --build
You will see NATS JetStream initialize, Redis bind to port 6379, and both Go microservices report Ready.🧪 Testing & VerificationOnce the containers are running, open a new terminal window to test the API and verify queue workers.Scenario A: Successful Async NotificationDispatch a standard email request to the API Gateway:Bashcurl -X POST http://localhost:8080/api/v1/notify \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "recruiter_test_1",
    "type": "EMAIL",
    "target": "evaluation@example.com",
    "title": "Interview Invitation",
    "body": "Testing the Go distributed worker pipeline."
  }'
Expected Response (Instant 202 Accepted):JSON{
  "notification_id": "8f9d2a1b-4c5e-6a7b-8c9d-0e1f2a3b4c5d",
  "status": "queued"
}
Scenario B: Testing the Dead Letter Queue (DLQ) & Circuit BreakerTo witness the self-healing retry loop in action, pass any target email containing the word "fail" while running the mock provider. This intentionally triggers a simulated third-party network outage:Bashcurl -X POST http://localhost:8080/api/v1/notify \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "recruiter_test_2",
    "type": "EMAIL",
    "target": "fail@example.com",
    "title": "System Alert",
    "body": "This payload will trigger the DLQ safety mechanism."
  }'
Watch the Worker Container Logs:Bashdocker compose logs -f worker
You will observe the worker automatically catch the failure, re-queue the job via negative acknowledgement (msg.Nak()), increment the delivery attempt counter, and safely divert the message to the DLQ stream:Plaintextworker-1  | [START] Processing EMAIL for fail@example.com (Attempt 1/3)
worker-1  | [MOCK ERROR] Simulated delivery rejection for: fail@example.com
worker-1  | [FAILED] Delivery failed... Re-queuing...
worker-1  | [START] Processing EMAIL for fail@example.com (Attempt 2/3)
worker-1  | [START] Processing EMAIL for fail@example.com (Attempt 3/3)
worker-1  | [DLQ] Message 8f9d2a1b... failed 3 times. Routing to DLQ.NOTIFY.
📁 Repository StructurePlaintext├── cmd/
│   ├── api/main.go          # HTTP Gateway entrypoint & REST handlers
│   └── worker/main.go       # NATS JetStream consumer & DLQ logic
├── internal/
│   ├── broker/nats.go       # NATS stream initialization & connection pooling
│   ├── cache/redis.go       # Redis client & idempotency rate-limiting
│   ├── models/payload.go    # Data structures & JSON serialization schemas
│   └── providers/           # Interface-driven delivery implementations
│       ├── provider.go      # NotificationProvider interface definition
│       ├── mock.go          # Mock logging provider for CI/CD & DLQ testing
│       ├── smtp.go          # Standard Net/SMTP integration (Gmail)
│       ├── sendgrid.go      # SendGrid REST SDK implementation
│       └── twilio.go        # Twilio HTTP API SMS client
├── Dockerfile               # Multi-stage Alpine build specs for tiny Go binaries
├── docker-compose.yml       # Local container orchestration & port mapping
└── .env.example             # Template for safe configuration management
