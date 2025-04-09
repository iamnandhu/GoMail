package auth

// LoginRequest represents a request to login a user
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a response from a login request
type LoginResponse struct {
	Token   string `json:"token"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// RegisterRequest represents a request to register a new user
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// RegisterResponse represents a response from a register request
type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// VerifyTokenRequest represents a request to verify a token
type VerifyTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// VerifyTokenResponse represents a response from a token verification
type VerifyTokenResponse struct {
	Valid bool   `json:"valid"`
	Error string `json:"error,omitempty"`
}

// LogoutRequest represents a request to logout a user
type LogoutRequest struct {
	Token string `json:"token" binding:"required"`
}

// LogoutResponse represents a response from a logout request
type LogoutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
} 