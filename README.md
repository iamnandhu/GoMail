# GoMail

GoMail is a lightweight, scalable SMTP client API built in Go. It provides a simple interface for sending emails and tracking email history.

## Project Structure

```
GoMail/
├── main.go                 # Application entry point
├── src/
│   ├── config/             # Configuration loading
│   ├── server/             # HTTP server setup
│   ├── middleware/         # HTTP middleware
│   ├── handler/            # HTTP request handlers
│   ├── logic/              # Business logic
│   ├── repository/         # Data access layer
│   ├── libs/               # Utility libraries
│   │   └── smtp/           # SMTP client implementation
│   └── utils/              # Helper utilities
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- MongoDB (for storing email history)

### Installation

1. Clone the repository
```bash
git clone https://github.com/username/GoMail.git
cd GoMail
```

2. Install dependencies
```bash
go mod download
```

3. Set up environment variables
```bash
export PORT=8080
export MONGODB_URI=mongodb://localhost:27017
export MONGODB_DB=gomail
export MONGODB_COLLECTION=emails
export SMTP_HOST=smtp.example.com
export SMTP_PORT=587
export SMTP_USERNAME=your-username
export SMTP_PASSWORD=your-password
export SMTP_FROM=no-reply@example.com
```

4. Run the application
```bash
go run main.go
```

## API Endpoints

- `POST /api/v1/email/send` - Send an email

## License

This project is licensed under the MIT License - see the LICENSE file for details.
