package repository

import (
	"GoMail/app/repository/email"
	"GoMail/app/repository/emaillog"
	"GoMail/app/repository/models"
	"context"

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
}

type repoImpl struct {
	email    email.EmailRepository
	emailLog emaillog.EmailLogRepository
}

func New(db *DB) Repository {
	return &repoImpl{
		email:    email.New(db.MongoDB),
		emailLog: emaillog.New(db.MongoDB),
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