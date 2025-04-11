package auth

import (
	"GoMail/app/repository/models"
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Register creates a new user account
func (s *service) Register(ctx context.Context, user *models.User, password string) error {
	// Check if user already exists
	existingUser, err := s.repo.FindUserByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	
	if existingUser != nil {
		return ErrUserAlreadyExists
	}
	
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	
	// Save user
	return s.repo.SaveUser(ctx, user)
} 