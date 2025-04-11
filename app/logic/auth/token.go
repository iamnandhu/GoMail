package auth

import (
	"GoMail/app/repository/models"
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// VerifyToken validates and parses a JWT token
func (s *service) VerifyToken(tokenString string) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.Secret), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	// Validate token
	if !token.Valid {
		return nil, ErrInvalidToken
	}
	
	// Extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}
	
	// Check if token revocation is enabled
	if s.config.JWT.EnableTokenRevoking && s.repo != nil {
		// Construct a context since we're not in a HTTP handler
		ctx := context.Background()
		
		// Check if token is revoked using token ID (jti claim)
		tokenID := claims.ID
		if tokenID == "" {
			// If no token ID, use the hash of the token string as ID
			tokenID = tokenString
		}
		
		isRevoked, err := s.repo.IsTokenRevoked(ctx, tokenID)
		if err != nil {
			return nil, err
		}
		
		if isRevoked {
			return nil, ErrInvalidToken
		}
	}
	
	return claims, nil
}

// generateToken generates a JWT token for a user
func (s *service) generateToken(user *models.User) (string, error) {
	// Set expiration time
	expirationTime := time.Now().Add(s.config.JWT.ExpiresIn)
	
	// Create claims
	claims := &Claims{
		UserID: user.ID.Hex(),
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.Hex(),
			ID:        generateTokenID(user.ID.Hex(), time.Now().Unix()),
		},
	}
	
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Sign token
	return token.SignedString([]byte(s.config.JWT.Secret))
}

// RevokeToken revokes a JWT token by adding it to the revoked tokens list
func (s *service) RevokeToken(ctx context.Context, tokenString string) error {
	// Skip if token revocation is disabled or repo is nil
	if !s.config.JWT.EnableTokenRevoking || s.repo == nil {
		return nil
	}
	
	// Parse token to get claims
	claims, err := s.VerifyToken(tokenString)
	if err != nil {
		return err
	}
	
	// Extract token ID
	tokenID := claims.ID
	if tokenID == "" {
		// If no token ID, use the hash of the token string as ID
		tokenID = tokenString
	}
	
	// Extract expiration time
	var expiresAt time.Time
	if claims.ExpiresAt != nil {
		expiresAt = claims.ExpiresAt.Time
	} else {
		// Use default expiration time
		expiresAt = time.Now().Add(s.config.JWT.ExpiresIn)
	}
	
	// Create revoked token record
	revokedToken := &models.RevokedToken{
		TokenID:   tokenID,
		UserID:    claims.UserID,
		RevokedAt: time.Now(),
		ExpiresAt: expiresAt,
		Reason:    "user_logout",
	}
	
	// Add to revoked tokens list
	return s.repo.RevokeToken(ctx, revokedToken)
}

// generateTokenID creates a unique ID for the token
func generateTokenID(userID string, timestamp int64) string {
	return userID + "-" + time.Unix(timestamp, 0).Format(time.RFC3339)
} 