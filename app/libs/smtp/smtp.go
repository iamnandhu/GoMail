package smtp

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"sync"
	"time"
)

// SMTPClient defines the interface for SMTP operations
type SMTPClient interface {
	Connect() error
	Disconnect() error
	Send(ctx context.Context, from, to, subject, body string) error
	SendHTML(ctx context.Context, from, to, subject, htmlBody string) error
	SendWithAttachments(ctx context.Context, from, to, subject, body string, attachments []Attachment) error
	SendBulk(ctx context.Context, messages []EmailRequest) []EmailResponse
	IsConnected() bool
}

// NewClient creates a new SMTP client
func NewClient(config Config) SMTPClient {
	return &smtpClient{
		config:        config,
		clientPool:    make([]*smtp.Client, 0, config.PoolSize),
		poolMutex:     &sync.Mutex{},
		retryAttempts: config.RetryAttempts,
		retryDelay:    config.RetryDelay,
	}
}

// smtpClient implements the Client interface
type smtpClient struct {
	config        Config
	client        *smtp.Client // Used for single connections
	clientPool    []*smtp.Client
	poolMutex     *sync.Mutex
	retryAttempts int
	retryDelay    time.Duration
}

// IsConnected checks if the client is connected
func (c *smtpClient) IsConnected() bool {
	if len(c.clientPool) > 0 {
		return true
	}
	return c.client != nil
}

// Connect establishes a connection to the SMTP server
func (c *smtpClient) Connect() error {
	// If already connected, return
	if c.IsConnected() {
		return nil
	}

	// If pooling is enabled, initialize the pool
	if c.config.PoolSize > 0 {
		return c.initializePool()
	}

	// Single connection mode
	return c.connectSingle()
}

// initializePool creates a pool of SMTP connections
func (c *smtpClient) initializePool() error {
	c.poolMutex.Lock()
	defer c.poolMutex.Unlock()

	// Clear any existing connections
	for _, client := range c.clientPool {
		client.Close()
	}
	c.clientPool = make([]*smtp.Client, 0, c.config.PoolSize)

	// Create new connections
	for i := 0; i < c.config.PoolSize; i++ {
		client, err := c.createConnection()
		if err != nil {
			// Close already created connections
			for _, existingClient := range c.clientPool {
				existingClient.Close()
			}
			c.clientPool = nil
			return fmt.Errorf("failed to initialize connection pool, connection %d failed: %w", i, err)
		}
		c.clientPool = append(c.clientPool, client)
	}

	return nil
}

// getClientFromPool gets an available client from the pool
func (c *smtpClient) getClientFromPool() (*smtp.Client, error) {
	c.poolMutex.Lock()
	defer c.poolMutex.Unlock()

	if len(c.clientPool) == 0 {
		return nil, errors.New("no connections available in the pool")
	}

	// Get a client from the pool
	client := c.clientPool[0]
	c.clientPool = c.clientPool[1:]

	// Check if the connection is still alive
	if err := client.Noop(); err != nil {
		// Connection is dead, create a new one
		newClient, err := c.createConnection()
		if err != nil {
			return nil, fmt.Errorf("failed to create new connection: %w", err)
		}
		client = newClient
	}

	return client, nil
}

// returnClientToPool returns a client to the pool
func (c *smtpClient) returnClientToPool(client *smtp.Client) {
	c.poolMutex.Lock()
	defer c.poolMutex.Unlock()

	// If the pool is full, close the client
	if len(c.clientPool) >= c.config.PoolSize {
		client.Close()
		return
	}

	c.clientPool = append(c.clientPool, client)
}

// connectSingle establishes a single connection
func (c *smtpClient) connectSingle() error {
	client, err := c.createConnection()
	if err != nil {
		return err
	}
	c.client = client
	return nil
}

// createConnection creates a new SMTP connection
func (c *smtpClient) createConnection() (*smtp.Client, error) {
	// Format server address
	addr := fmt.Sprintf("%s:%s", c.config.Host, c.config.Port)

	// Connect to the SMTP server
	var client *smtp.Client
	var err error

	// Set up dialer with timeout
	dialer := &net.Dialer{
		Timeout: c.config.ConnectTimeout,
	}

	if c.config.UseTLS {
		// Connect with TLS
		tlsConfig := &tls.Config{
			ServerName:         c.config.Host,
			InsecureSkipVerify: c.config.InsecureSkipVerify,
		}
		conn, err := tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to SMTP server with TLS: %w", err)
		}
		client, err = smtp.NewClient(conn, c.config.Host)
	} else {
		// Connect without TLS
		conn, err := dialer.Dial("tcp", addr)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to SMTP server: %w", err)
		}
		client, err = smtp.NewClient(conn, c.config.Host)

		// Start TLS if required
		if c.config.StartTLS {
			tlsConfig := &tls.Config{
				ServerName:         c.config.Host,
				InsecureSkipVerify: c.config.InsecureSkipVerify,
			}
			if err = client.StartTLS(tlsConfig); err != nil {
				client.Close()
				return nil, fmt.Errorf("failed to start TLS: %w", err)
			}
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to SMTP server: %w", err)
	}

	// Authenticate if credentials are provided
	if c.config.Username != "" && c.config.Password != "" {
		auth := smtp.PlainAuth("", c.config.Username, c.config.Password, c.config.Host)
		if err := client.Auth(auth); err != nil {
			client.Close()
			return nil, fmt.Errorf("SMTP authentication failed: %w", err)
		}
	}

	return client, nil
}

// Disconnect closes the connection(s) to the SMTP server
func (c *smtpClient) Disconnect() error {
	var lastErr error

	// Close the connection pool if it exists
	if len(c.clientPool) > 0 {
		c.poolMutex.Lock()
		for _, client := range c.clientPool {
			if err := client.Quit(); err != nil {
				lastErr = fmt.Errorf("failed to disconnect from SMTP server: %w", err)
			}
		}
		c.clientPool = nil
		c.poolMutex.Unlock()
	}

	// Close the single client if it exists
	if c.client != nil {
		err := c.client.Quit()
		c.client = nil
		if err != nil {
			lastErr = fmt.Errorf("failed to disconnect from SMTP server: %w", err)
		}
	}

	return lastErr
}

// Send sends a plain text email through the SMTP server
func (c *smtpClient) Send(ctx context.Context, from, to, subject, body string) error {
	req := EmailRequest{
		From:    from,
		To:      to,
		Subject: subject,
		Body:    body,
		IsHTML:  false,
	}
	return c.sendWithRetry(ctx, req)
}

// SendHTML sends an HTML email through the SMTP server
func (c *smtpClient) SendHTML(ctx context.Context, from, to, subject, htmlBody string) error {
	req := EmailRequest{
		From:    from,
		To:      to,
		Subject: subject,
		Body:    htmlBody,
		IsHTML:  true,
	}
	return c.sendWithRetry(ctx, req)
}

// SendWithAttachments sends an email with attachments through the SMTP server
func (c *smtpClient) SendWithAttachments(ctx context.Context, from, to, subject, body string, attachments []Attachment) error {
	req := EmailRequest{
		From:        from,
		To:          to,
		Subject:     subject,
		Body:        body,
		IsHTML:      false,
		Attachments: attachments,
	}
	return c.sendWithRetry(ctx, req)
}

// SendBulk sends multiple emails concurrently
func (c *smtpClient) SendBulk(ctx context.Context, messages []EmailRequest) []EmailResponse {
	responses := make([]EmailResponse, len(messages))
	var wg sync.WaitGroup

	// Use a semaphore to limit concurrent goroutines
	semaphore := make(chan struct{}, c.config.MaxConcurrent)

	for i, msg := range messages {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire semaphore

		go func(idx int, email EmailRequest) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore

			var err error
			// Use retry mechanism
			for attempt := 0; attempt <= c.retryAttempts; attempt++ {
				select {
				case <-ctx.Done():
					responses[idx] = EmailResponse{
						Success: false,
						Error:   "context cancelled",
					}
					return
				default:
					if attempt > 0 && c.retryDelay > 0 {
						time.Sleep(c.retryDelay)
					}

					err = c.sendEmail(ctx, email)
					if err == nil {
						responses[idx] = EmailResponse{Success: true}
						return
					}
				}
			}

			// All attempts failed
			responses[idx] = EmailResponse{
				Success: false,
				Error:   err.Error(),
			}
		}(i, msg)
	}

	wg.Wait()
	return responses
}

// sendWithRetry attempts to send an email with retries
func (c *smtpClient) sendWithRetry(ctx context.Context, req EmailRequest) error {
	var lastErr error

	for attempt := 0; attempt <= c.retryAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if attempt > 0 && c.retryDelay > 0 {
				time.Sleep(c.retryDelay)
			}

			err := c.sendEmail(ctx, req)
			if err == nil {
				return nil
			}
			lastErr = err
		}
	}

	return fmt.Errorf("all send attempts failed, last error: %w", lastErr)
}

// sendEmail sends a single email
func (c *smtpClient) sendEmail(ctx context.Context, req EmailRequest) error {
	// Connect if not already connected
	if !c.IsConnected() {
		if err := c.Connect(); err != nil {
			return err
		}
	}

	// Get a client to use (either from pool or the single client)
	var client *smtp.Client
	var err error
	returnToPool := false

	if c.config.PoolSize > 0 {
		client, err = c.getClientFromPool()
		if err != nil {
			return err
		}
		returnToPool = true
	} else {
		client = c.client
	}

	// Ensure we return the client to the pool
	defer func() {
		if returnToPool {
			c.returnClientToPool(client)
		}
	}()

	// Use default sender if not specified
	from := req.From
	if from == "" {
		from = c.config.From
	}

	// Prepare email headers and body
	var message string
	if req.IsHTML {
		message = buildHTMLEmail(from, req.To, req.Subject, req.Body)
	} else if len(req.Attachments) > 0 {
		message = buildMultipartEmail(from, req.To, req.Subject, req.Body, req.Attachments)
	} else {
		message = buildPlainEmail(from, req.To, req.Subject, req.Body)
	}

	// Set the sender
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set the recipients
	recipients := strings.Split(req.To, ",")
	for _, recipient := range recipients {
		recipient = strings.TrimSpace(recipient)
		if recipient != "" {
			if err := client.Rcpt(recipient); err != nil {
				return fmt.Errorf("failed to add recipient %s: %w", recipient, err)
			}
		}
	}

	// Send the email body
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data writer: %w", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write email data: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	return nil
}

// Helper functions to build email content
func buildPlainEmail(from, to, subject, body string) string {
	return fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"Content-Type: text/plain; charset=UTF-8\r\n\r\n"+
			"%s",
		from, to, subject, body)
}

func buildHTMLEmail(from, to, subject, htmlBody string) string {
	return fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n\r\n"+
			"%s",
		from, to, subject, htmlBody)
}

func buildMultipartEmail(from, to, subject, body string, attachments []Attachment) string {
	// Generate a boundary string
	boundary := fmt.Sprintf("_boundary_%d", time.Now().UnixNano())

	// Build the MIME multipart message
	var builder strings.Builder

	// Add headers
	builder.WriteString(fmt.Sprintf("From: %s\r\n", from))
	builder.WriteString(fmt.Sprintf("To: %s\r\n", to))
	builder.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	builder.WriteString("MIME-Version: 1.0\r\n")
	builder.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n\r\n", boundary))

	// Add text part
	builder.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	builder.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
	builder.WriteString(body)
	builder.WriteString("\r\n\r\n")

	// Add attachments
	for _, att := range attachments {
		builder.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		builder.WriteString(fmt.Sprintf("Content-Type: %s\r\n", att.MimeType))
		builder.WriteString("Content-Transfer-Encoding: base64\r\n")
		builder.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\r\n\r\n", att.Filename))

		// TODO: Convert attachment content to base64
		// This is a simplified version; you'd need to add base64 encoding
		builder.WriteString("[BASE64_ENCODED_CONTENT]")
		builder.WriteString("\r\n\r\n")
	}

	// Close the MIME multipart message
	builder.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	return builder.String()
}