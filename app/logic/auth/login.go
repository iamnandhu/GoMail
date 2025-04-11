package auth

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

// Login authenticates a user and returns a JWT token
func (s *service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	
	if user == nil {
		return "", ErrInvalidCredentials
	}
	
	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}
	
	// Generate JWT token
	return s.generateToken(user)
} 