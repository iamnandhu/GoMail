package smtp

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
	"time"
)

// SMTPClient defines the interface for SMTP operations
type SMTPClient interface {
	Connect() error
	Disconnect() error
	Send(ctx context.Context, from, to, subject, body string) error
	SendHTML(ctx context.Context, from, to, subject, htmlBody string) error
	SendWithAttachments(ctx context.Context, from, to, subject, body string, attachments []Attachment) error
	IsConnected() bool
}

// NewClient creates a new SMTP client
func NewClient(config Config) SMTPClient {
	return &smtpClient{
		config:        config,
		clientPool:    make(chan *smtp.Client, config.PoolSize),
		retryAttempts: config.RetryAttempts,
		retryDelay:    config.RetryDelay,
	}
}

// smtpClient implements the Client interface
type smtpClient struct {
	config        Config
	client        *smtp.Client // Used for single connections
	clientPool    chan *smtp.Client
	retryAttempts int
	retryDelay    time.Duration
}

// IsConnected checks if the client is connected
func (c *smtpClient) IsConnected() bool {
	if c.clientPool != nil && cap(c.clientPool) > 0 && len(c.clientPool) > 0 {
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
	// Clear existing pool if it exists
	if c.clientPool != nil {
		close(c.clientPool)
	}
	
	// Create a new pool
	c.clientPool = make(chan *smtp.Client, c.config.PoolSize)

	// Create new connections
	for i := 0; i < c.config.PoolSize; i++ {
		client, err := c.createConnection()
		if err != nil {
			// Close the pool and any clients already created
			for i := 0; i < len(c.clientPool); i++ {
				select {
				case client := <-c.clientPool:
					client.Close()
				default:
					break
				}
			}
			close(c.clientPool)
			c.clientPool = nil
			return fmt.Errorf("failed to initialize connection pool, connection %d failed: %w", i, err)
		}
		c.clientPool <- client
	}

	return nil
}

// getClientFromPool gets an available client from the pool
func (c *smtpClient) getClientFromPool() (*smtp.Client, error) {
	// Get a client from the pool
	var client *smtp.Client
	select {
	case client = <-c.clientPool:
		// Check if the connection is still alive
		if err := client.Noop(); err != nil {
			// Connection is dead, create a new one
			log.Printf("SMTP connection health check failed: %v, creating new connection", err)
			newClient, err := c.createConnection()
			if err != nil {
				return nil, fmt.Errorf("failed to create new connection: %w", err)
			}
			client = newClient
		}
		return client, nil
	default:
		return nil, errors.New("no connections available in the pool")
	}
}

// returnClientToPool returns a client to the pool
func (c *smtpClient) returnClientToPool(client *smtp.Client) {
	// Try to return to the pool, but don't block if full
	select {
	case c.clientPool <- client:
		// Successfully returned to pool
	default:
		// Pool is full, close the client
		client.Close()
	}
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
	
	// Debug: Log SMTP configuration
	log.Printf("SMTP Config - Host: %s, Port: %s", c.config.Host, c.config.Port)

	// Connect to the SMTP server
	var client *smtp.Client
	var err error

	// Set up dialer with timeout
	dialer := &net.Dialer{
		Timeout: c.config.ConnectTimeout,
	}

	if c.config.UseTLS {
		// Debug: Log TLS connection info
		log.Printf("Connecting with TLS to %s", addr)
		
		// Connect with TLS
		tlsConfig := &tls.Config{
			ServerName:         c.config.Host,
			InsecureSkipVerify: c.config.InsecureSkipVerify,
		}
		conn, err := tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
		if err != nil {
			log.Printf("TLS connection error: %v", err)
			return nil, fmt.Errorf("failed to connect to SMTP server with TLS: %w", err)
		}
		client, err = smtp.NewClient(conn, c.config.Host)
	} else {
		// Debug: Log non-TLS connection info
		log.Printf("Connecting without TLS to %s", addr)
		
		// Connect without TLS
		conn, err := dialer.Dial("tcp", addr)
		if err != nil {
			log.Printf("Connection error: %v", err)
			return nil, fmt.Errorf("failed to connect to SMTP server: %w", err)
		}
		client, err = smtp.NewClient(conn, c.config.Host)

		// Start TLS if required
		if c.config.StartTLS {
			log.Printf("Starting TLS after connection")
			tlsConfig := &tls.Config{
				ServerName:         c.config.Host,
				InsecureSkipVerify: c.config.InsecureSkipVerify,
			}
			if err = client.StartTLS(tlsConfig); err != nil {
				log.Printf("StartTLS error: %v", err)
				client.Close()
				return nil, fmt.Errorf("failed to start TLS: %w", err)
			}
		}
	}

	if err != nil {
		log.Printf("Client creation error: %v", err)
		return nil, fmt.Errorf("failed to connect to SMTP server: %w", err)
	}

	// Authenticate if credentials are provided
	if c.config.Username != "" && c.config.Password != "" {
		log.Printf("Authenticating with username: %s", c.config.Username)
		auth := smtp.PlainAuth("", c.config.Username, c.config.Password, c.config.Host)
		if err := client.Auth(auth); err != nil {
			log.Printf("Authentication error: %v", err)
			client.Close()
			return nil, fmt.Errorf("SMTP authentication failed: %w", err)
		}
		log.Printf("Authentication successful")
	}

	log.Printf("SMTP connection established successfully")
	return client, nil
}

// Disconnect closes the connection(s) to the SMTP server
func (c *smtpClient) Disconnect() error {
	var lastErr error

	// Close the connection pool if it exists
	if c.clientPool != nil && cap(c.clientPool) > 0 {
		// Drain and close all clients in the pool
		for {
			select {
			case client := <-c.clientPool:
				if err := client.Quit(); err != nil {
					lastErr = fmt.Errorf("failed to disconnect from SMTP server: %w", err)
				}
			default:
				// Pool is empty
				close(c.clientPool)
				c.clientPool = nil
				goto CleanupSingleClient
			}
		}
	}

CleanupSingleClient:
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
			log.Printf("Email send attempt %d failed: %v", attempt+1, err)
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
	fromAddr := parseAddress(from)
	toAddr := parseAddress(to)
	
	header := make(map[string]string)
	header["From"] = fromAddr
	header["To"] = toAddr
	header["Subject"] = subject
	header["Content-Type"] = "text/plain; charset=UTF-8"
	
	return buildMessage(header, body)
}

func buildHTMLEmail(from, to, subject, htmlBody string) string {
	fromAddr := parseAddress(from)
	toAddr := parseAddress(to)
	
	header := make(map[string]string)
	header["From"] = fromAddr
	header["To"] = toAddr
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=UTF-8"
	
	return buildMessage(header, htmlBody)
}

func buildMultipartEmail(from, to, subject, body string, attachments []Attachment) string {
	// Generate a boundary string
	boundary := fmt.Sprintf("_boundary_%d", time.Now().UnixNano())

	// Create buffer for multipart message
	var buf strings.Builder
	
	// Set up headers
	fromAddr := parseAddress(from)
	toAddr := parseAddress(to)
	
	// Write the headers
	buf.WriteString(fmt.Sprintf("From: %s\r\n", fromAddr))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", toAddr))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n\r\n", boundary))

	// Add text part
	buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	buf.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
	buf.WriteString(body)
	buf.WriteString("\r\n\r\n")

	// Add attachments
	for _, att := range attachments {
		buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		buf.WriteString(fmt.Sprintf("Content-Type: %s\r\n", att.MimeType))
		buf.WriteString("Content-Transfer-Encoding: base64\r\n")
		buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\r\n\r\n", att.Filename))

		// Base64 encode the attachment content
		encoded := base64.StdEncoding.EncodeToString(att.Content)
		
		// Write in lines of 76 characters as per RFC 2045
		for i := 0; i < len(encoded); i += 76 {
			end := i + 76
			if end > len(encoded) {
				end = len(encoded)
			}
			buf.WriteString(encoded[i:end])
			buf.WriteString("\r\n")
		}
		buf.WriteString("\r\n")
	}

	// Close the MIME multipart message
	buf.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	return buf.String()
}

// Helper function to build a message with headers
func buildMessage(header map[string]string, body string) string {
	var buf strings.Builder
	
	// Add headers
	for key, value := range header {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	
	// Add separator between headers and body
	buf.WriteString("\r\n")
	
	// Add body
	buf.WriteString(body)
	
	return buf.String()
}

// Helper function to parse email addresses
func parseAddress(addr string) string {
	// If already properly formatted or empty, return as is
	if addr == "" || strings.HasPrefix(addr, "<") && strings.HasSuffix(addr, ">") {
		return addr
	}
	
	// Parse and format the email address
	if a, err := mail.ParseAddress(addr); err == nil {
		return a.String()
	}
	
	// Fallback to original address if parsing fails
	return addr
}