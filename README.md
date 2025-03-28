# Log-Flow: Real-Time Log Processing API

A high-performance, real-time log processing backend service built with Go, featuring asynchronous processing using RabbitMQ and WebSocket-based live updates. This RESTful API service provides endpoints for log file processing, real-time progress tracking, and comprehensive statistics.

## 🚀 Features

- **Asynchronous Log Processing**: Utilizes RabbitMQ for reliable message queuing and processing
- **Real-time Updates**: WebSocket integration for live progress tracking
- **Authentication**: Secure endpoints with middleware
- **Rate Limiting**: Protected API endpoints with configurable rate limits
- **Concurrent Processing**: Efficient log processing with Go routines
- **Storage Integration**: Flexible storage interface for log files
- **REST API**: Comprehensive API endpoints for log management
- **Clean Architecture**: Follows SOLID principles with clear separation of concerns and dependency injection
- **Fault Tolerance**: Robust retry mechanism with failed queue for manual inspection
- **Live Progress Tracking**: Real-time processing status via WebSocket connection

## 🛠 Tech Stack

- **Backend**: Go (Fiber framework)
- **Message Queue**: RabbitMQ
- **WebSocket**: Fiber WebSocket
- **Database & Auth**: 
  - Supabase for PostgreSQL database
  - Supabase Authentication with JWT
  - Supabase Storage for log file management
- **Rate Limiting**: Built-in rate limiting middleware

## 📋 API Endpoints

### Authentication Routes
```
POST /auth/login    - User login
POST /auth/register - User registration
```

### Log Management Routes
```
POST /api/upload-logs           - Upload log files for processing
GET  /api/stats                - Fetch aggregated statistics
GET  /api/stats/:jobId         - Fetch statistics for specific job
GET  /api/queue-status         - Get current queue status
GET  /api/live-stats/:jobID    - WebSocket endpoint for real-time updates
```

## 🔒 Security

- JWT-based authentication(Supabase Auth)
- Rate limiting on sensitive endpoints(Taking X-Real-IP if available via proxies like nginx, to prevent DOS attack using IP spoofing)
- Secure WebSocket connections
- Job-level authorization checks(Job infos can be accessed only by respective owners)

## 🎯 Performance

The service is designed to handle large log files efficiently:
- Concurrent processing using Go routines
- Streaming file processing to manage memory usage
- RabbitMQ for reliable message queuing
- WebSocket for efficient real-time updates

## 🔄 Real-Time Progress Tracking

The application features a WebSocket-based real-time progress tracking system:

- **Secure WebSocket Endpoint**: `/api/live-stats/:jobID` with job-level authorization
- **Live Updates Structure**:
  ```json
  {
    "jobID": "ae316344-3fe9-4770-b9a8-3114154165d9",
    "progressInPercentage": 20,
    "uniqueIPs": 1,
    "invalidLogs": 0,
    "totalLogsProcessed": 4,
    "status": "In Progress",
    "logLevelCounts": {
        "error": 0,
        "info": 2,
        "warn": 1
    },
    "keyWordCounts": {
        "error": 5,
        "exception": 2,
        "failed": 1
    }
  }
  ```
  
- **Real-Time Metrics**:
  - Processing Progress: Overall completion percentage
  - Unique IP Addresses: Count of distinct IPs found
  - Log Quality: Track of valid vs invalid log entries
  - Processing Volume: Total number of logs processed
  - Job Status: Current status of the processing job
  - Log Level Distribution: Counts by log level (error, info, warn, etc.)
  - Keyword Tracking: Frequency count of configured keywords

## 🔁 Fault Tolerance & Recovery

The system implements a robust error handling and recovery mechanism:

- **Retry Mechanism**:
  - Maximum 3 retry attempts for failed jobs
  - Time gap between retry attempts
  - Automatic tracking of retry counts
  
- **Failed Queue System**:
  - Failed jobs (after 3 retries) are moved to a dedicated `failed_queue`
  - Enables manual inspection and debugging
  - Provides ability to reprocess failed jobs after fixing issues
  - Maintains full error context and processing history

## 🚀 Getting Started

### Prerequisites

- Go 1.23+
- RabbitMQ
- PostgreSQL (or your preferred database)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/AbdulRahimOM/Log-Flow.git
cd log-flow
```

2. Install dependencies:
```bash
go mod download
```

3. Start RabbitMQ 
(using Docker):
```bash
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management
```

4. Run the application:
```bash
go run cmd/main.go
```

## 🏗 Project Structure

```
.
├── cmd/
│   ├── api/                    
│   │   └── main.go            # Main application Entry point
│   └── migrate/               
│       └── migrate.go         # Database migrations
├── internal/
│   ├── api/
│   │   ├── handler/           # Request handlers (auth, logs, websocket)
│   │   ├── middleware/        # Middlewares (auth, authorization, rate limiting)
│   │   └── routes/           # Route definitions
│   ├── domain/
│   │   ├── models/           # Data models and entities
│   │   └── response/         # Response structures and interfaces
│   ├── infrastructure/
│   │   ├── config/           # Configuration management
│   │   ├── db/              # Database connection and setup
│   │   ├── queue/           # RabbitMQ implementation
│   │   ├── server/          # Server initialization
│   │   └── storage/         # Storage interfaces and implementations
│   ├── utils/
│   │   ├── helper/          # Helper functions
│   │   ├── jwt/            # JWT implementation
│   │   ├── locals/         # Context utilities
│   │   └── validation/     # Request validation
│   └── workers/            # Worker implementations
├── Dockerfile             # Container configuration
├── docker-compose.yml     # Docker services orchestration
└── Makefile              # Build and development commands
```
