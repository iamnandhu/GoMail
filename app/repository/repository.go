package repository

import (
	"GoMail/app/repository/email"
	"GoMail/app/repository/emaillog"
	"GoMail/app/repository/models"
	"GoMail/app/repository/token"
	"GoMail/app/repository/user"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type DB struct {
	MongoDB *mongo.Database
}

type Repository interface {
	SaveEmail(ctx context.Context, email *models.Email) error
	FindEmails(ctx context.Context, filter interface{}, page, limit int) ([]*models.Email, int64, error)
	DeleteEmail(ctx context.Context, id interface{}) error
	
	// Email Log methods
	SaveEmailLog(ctx context.Context, emailLog *models.EmailLog) error
	FindEmailLogs(ctx context.Context, filter interface{}, page, limit int) ([]*models.EmailLog, int64, error)
	FindEmailLogByID(ctx context.Context, id string) (*models.EmailLog, error)
	
	// User methods
	SaveUser(ctx context.Context, user *models.User) error
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	FindUserByID(ctx context.Context, id string) (*models.User, error)
	FindUsers(ctx context.Context, filter interface{}, page, limit int) ([]*models.User, int64, error)
	
	// Token methods
	RevokeToken(ctx context.Context, token *models.RevokedToken) error
	IsTokenRevoked(ctx context.Context, tokenID string) (bool, error)
	CleanupExpiredTokens(ctx context.Context, beforeTime time.Time) (int64, error)
	InitTokenIndexes(ctx context.Context) error
}

type repoImpl struct {
	email    email.EmailRepository
	emailLog emaillog.EmailLogRepository
	user     user.Repository
	token    token.Repository
}

func New(db *DB) Repository {
	return &repoImpl{
		email:    email.New(db.MongoDB),
		emailLog: emaillog.New(db.MongoDB),
		user:     user.New(db.MongoDB),
		token:    token.New(db.MongoDB),
	}
}

func (r *repoImpl) SaveEmail(ctx context.Context, email *models.Email) error {
	return r.email.Save(ctx, email)
}

func (r *repoImpl) FindEmails(ctx context.Context, filter interface{}, page, limit int) ([]*models.Email, int64, error) {
	return r.email.FindAll(ctx, filter, page, limit)
}

func (r *repoImpl) DeleteEmail(ctx context.Context, id interface{}) error {
	return r.email.Delete(ctx, id)
}

// SaveEmailLog saves an email log to the database
func (r *repoImpl) SaveEmailLog(ctx context.Context, emailLog *models.EmailLog) error {
	return r.emailLog.Save(ctx, emailLog)
}

// FindEmailLogs retrieves email logs from the database
func (r *repoImpl) FindEmailLogs(ctx context.Context, filter interface{}, page, limit int) ([]*models.EmailLog, int64, error) {
	return r.emailLog.FindAll(ctx, filter, page, limit)
}

// FindEmailLogByID retrieves an email log by ID
func (r *repoImpl) FindEmailLogByID(ctx context.Context, id string) (*models.EmailLog, error) {
	return r.emailLog.FindByID(ctx, id)
}

// SaveUser saves a user to the database
func (r *repoImpl) SaveUser(ctx context.Context, user *models.User) error {
	return r.user.Save(ctx, user)
}

// FindUserByEmail finds a user by email
func (r *repoImpl) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return r.user.FindByEmail(ctx, email)
}

// FindUserByID finds a user by ID
func (r *repoImpl) FindUserByID(ctx context.Context, id string) (*models.User, error) {
	return r.user.FindByID(ctx, id)
}

// FindUsers retrieves users from the database
func (r *repoImpl) FindUsers(ctx context.Context, filter interface{}, page, limit int) ([]*models.User, int64, error) {
	return r.user.FindAll(ctx, filter, page, limit)
}

// RevokeToken adds a token to the revoked tokens list
func (r *repoImpl) RevokeToken(ctx context.Context, token *models.RevokedToken) error {
	return r.token.RevokeToken(ctx, token)
}

// IsTokenRevoked checks if a token is revoked
func (r *repoImpl) IsTokenRevoked(ctx context.Context, tokenID string) (bool, error) {
	return r.token.IsTokenRevoked(ctx, tokenID)
}

// CleanupExpiredTokens removes expired tokens
func (r *repoImpl) CleanupExpiredTokens(ctx context.Context, beforeTime time.Time) (int64, error) {
	return r.token.CleanupExpiredTokens(ctx, beforeTime)
}

// InitTokenIndexes initializes indexes for token collection
func (r *repoImpl) InitTokenIndexes(ctx context.Context) error {
	return r.token.CreateTokenRevokedIndex(ctx)
} 