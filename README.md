# 📧 GoMail

<div align="center">
  
![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-green?style=for-the-badge)
![Status](https://img.shields.io/badge/status-active-success?style=for-the-badge)

</div>

GoMail is a lightweight, high-performance SMTP client API built in Go. It provides a simple interface for sending emails with features like connection pooling, retries, and comprehensive email tracking.

## ✨ Features

- 🚀 **High Performance** - Connection pooling and concurrent email sending
- 🔄 **Retry Mechanism** - Automatic retries for failed email attempts
- 📊 **Email Tracking** - MongoDB integration for email history and analytics
- 🔒 **Authentication** - JWT-based authentication for API endpoints
- 📝 **Rich Content** - Support for HTML emails and attachments
- 🔄 **Bulk Operations** - Send multiple emails in a single request
- 🛠️ **Configurable** - Extensive configuration options via YAML or environment variables

## 🏗️ Project Structure

```
GoMail/
├── app/
│   ├── config/             # Configuration loading
│   ├── server/             # HTTP server setup
│   ├── middleware/         # HTTP middleware
│   ├── handler/            # HTTP request handlers
│   ├── logic/              # Business logic
│   ├── repository/         # Data access layer
│   ├── libs/               # Utility libraries
│   │   └── smtp/           # SMTP client implementation
│   └── utils/              # Helper utilities
├── main.go                 # Application entry point
├── config.yaml             # Configuration file
└── go.mod                  # Go module file
```

## 🚀 Getting Started

### Prerequisites

- Go 1.23+ 
- MongoDB (for storing email history)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/username/GoMail.git
   cd GoMail
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure the application**
   
   Create a `config.yaml` file in the root directory or set environment variables:
   
   ```yaml
   # config.yaml example
   env: "dev"
   
   server:
     port: "8080"
     portGRPC: "50051"
     readTimeout: 10s
     writeTimeout: 10s
   
   mongodb:
     uri: "mongodb://localhost:27017/"
     username: "username"
     password: "password123"
     database: "gomail"
     endpoint: "localhost:27017"
     timeout: 10s
     connectionTimeout: 10s
   
   smtp:
     host: "smtp.example.com"
     port: "587"
     username: "your-smtp-username"
     password: "your-smtp-password"
     from: "no-reply@example.com"
     useStartTLS: true
     maxConcurrent: 10
   
   jwt:
     secret: "your_jwt_secret_key_change_in_production"
     expiresIn: 24h
   ```
   
   Alternatively, you can use environment variables:
   
   ```bash
   export PORT=8080
   export MONGODB_URI=mongodb://localhost:27017
   export MONGODB_USERNAME=username
   export MONGODB_PASSWORD=password123
   export MONGODB_DATABASE=gomail
   export SMTP_HOST=smtp.example.com
   export SMTP_PORT=587
   export SMTP_USERNAME=your-username
   export SMTP_PASSWORD=your-password
   export SMTP_FROM=no-reply@example.com
   export JWT_SECRET=your_jwt_secret_key
   ```

4. **Run the application**
   ```bash
   go run main.go
   ```

## 🔌 API Reference

GoMail provides a RESTful API for sending emails. All API endpoints are prefixed with `/api/v1`.

### Authentication

Most endpoints require authentication. You need to include a JWT token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

You can obtain a token by using the login endpoint:

```bash
POST /api/v1/auth/login
```

**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "success": true
}
```

### Email Endpoints

#### Check Email Service Status

```bash
GET /api/v1/email/status
```

**Response:**
```json
{
  "status": "operational"
}
```

#### Send Plain Text Email

```bash
POST /api/v1/email/send
```

**Request:**
```json
{
  "from": "sender@example.com",
  "to": "recipient@example.com",
  "subject": "Hello from GoMail",
  "body": "This is a test email sent using GoMail API."
}
```

**Response:**
```json
{
  "success": true
}
```

#### Send HTML Email

```bash
POST /api/v1/email/send-html
```

**Request:**
```json
{
  "from": "sender@example.com",
  "to": "recipient@example.com",
  "subject": "Hello from GoMail",
  "body": "<h1>Hello!</h1><p>This is a <strong>HTML</strong> email sent using GoMail API.</p>"
}
```

**Response:**
```json
{
  "success": true
}
```

#### Send Email with Attachments

```bash
POST /api/v1/email/send-with-attachments
```

**Request:**
```json
{
  "from": "sender@example.com",
  "to": "recipient@example.com",
  "subject": "Email with Attachment",
  "body": "Please find the attached document.",
  "attachments": [
    {
      "filename": "document.pdf",
      "content": "base64-encoded-content",
      "mimeType": "application/pdf"
    }
  ]
}
```

**Response:**
```json
{
  "success": true
}
```

#### Send Bulk Emails

```bash
POST /api/v1/email/send-bulk
```

**Request:**
```json
{
  "emails": [
    {
      "from": "sender@example.com",
      "to": "recipient1@example.com",
      "subject": "Hello Recipient 1",
      "body": "This is a message for recipient 1",
      "isHtml": false
    },
    {
      "from": "sender@example.com",
      "to": "recipient2@example.com",
      "subject": "Hello Recipient 2",
      "body": "<h1>Hello Recipient 2</h1>",
      "isHtml": true
    }
  ]
}
```

**Response:**
```json
{
  "results": [
    {
      "success": true
    },
    {
      "success": true
    }
  ]
}
```

## ⚙️ Configuration Options

GoMail offers extensive configuration options through YAML or environment variables:

### Server Configuration

| YAML Key | Environment Variable | Description | Default |
|----------|----------------------|-------------|---------|
| `server.port` | `PORT` | HTTP server port | `8080` |
| `server.portGRPC` | `PORT_GRPC` | gRPC server port | `50051` |
| `server.readTimeout` | - | HTTP read timeout | `10s` |
| `server.writeTimeout` | - | HTTP write timeout | `10s` |

### MongoDB Configuration

| YAML Key | Environment Variable | Description | Default |
|----------|----------------------|-------------|---------|
| `mongodb.uri` | `MONGODB_URI` | MongoDB connection URI | - |
| `mongodb.username` | `MONGODB_USERNAME` | MongoDB username | - |
| `mongodb.password` | `MONGODB_PASSWORD` | MongoDB password | - |
| `mongodb.database` | `MONGODB_DATABASE` | MongoDB database name | `gomail` |
| `mongodb.endpoint` | `MONGODB_ENDPOINT` | MongoDB endpoint | `localhost:27017` |
| `mongodb.timeout` | - | MongoDB operation timeout | `10s` |
| `mongodb.connectionTimeout` | - | MongoDB connection timeout | `10s` |

### SMTP Configuration

| YAML Key | Environment Variable | Description | Default |
|----------|----------------------|-------------|---------|
| `smtp.host` | `SMTP_HOST` | SMTP server host | - |
| `smtp.port` | `SMTP_PORT` | SMTP server port | `587` |
| `smtp.username` | `SMTP_USERNAME` | SMTP username | - |
| `smtp.password` | `SMTP_PASSWORD` | SMTP password | - |
| `smtp.from` | `SMTP_FROM` | Default sender email | - |
| `smtp.useStartTLS` | `SMTP_USE_STARTTLS` | Use STARTTLS | `true` |
| `smtp.maxConcurrent` | `SMTP_MAX_CONCURRENT` | Max concurrent connections | `10` |

### JWT Configuration

| YAML Key | Environment Variable | Description | Default |
|----------|----------------------|-------------|---------|
| `jwt.secret` | `JWT_SECRET` | JWT secret key | - |
| `jwt.expiresIn` | - | JWT expiration time | `24h` |
| `jwt.enableTokenRevoking` | `JWT_ENABLE_TOKEN_REVOKING` | Enable token revocation | `false` |

## 📚 Advanced Usage

### Connection Pooling

GoMail supports SMTP connection pooling for improved performance. Configure this in your `config.yaml`:

```yaml
smtp:
  # ... other smtp settings
  maxConcurrent: 20  # Defines the connection pool size
```

### Retry Mechanism

Configure retry attempts for failed email sending:

```yaml
smtp:
  # ... other smtp settings
  retryAttempts: 3
  retryDelay: 5s
```

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgements

- [Go Gin Framework](https://github.com/gin-gonic/gin)
- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver)
- [Go YAML](https://github.com/go-yaml/yaml)
