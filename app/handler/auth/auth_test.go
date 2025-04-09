package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Test request bodies
var (
	validLoginRequestBody   = []byte(`{"email":"test@example.com","password":"password123"}`)
	invalidLoginRequestBody = []byte(`{"email":"","password":"password123"}`)
	
	validRegisterRequestBody   = []byte(`{"email":"newuser@example.com","password":"password123","firstName":"John","lastName":"Doe"}`)
	invalidRegisterRequestBody = []byte(`{"email":"invalidemail","password":"short"}`)
	
	validVerifyTokenRequestBody   = []byte(`{"token":"valid-token"}`)
	invalidVerifyTokenRequestBody = []byte(`{"token":""}`)
)

func Test_handler_login(t *testing.T) {
	type args struct {
		c       *gin.Context
		request []byte
	}
	tests := []struct {
		name               string
		args               args
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "happy path",
			args: args{
				c:       nil,
				request: validLoginRequestBody,
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: LoginResponse{
				Token:   "placeholder-token",
				Success: true,
			},
		},
		{
			name: "invalid request",
			args: args{
				c:       nil,
				request: invalidLoginRequestBody,
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Assemble
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			payload := bytes.NewBuffer(tt.args.request)
			r, _ := http.NewRequest("POST", "/auth/login", payload)
			r.Header.Set("Content-Type", "application/json")
			c.Request = r

			h := &Handler{}

			h.login(c)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			if tt.expectedStatusCode == http.StatusOK && tt.expectedResponse != nil {
				var response LoginResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.expectedResponse, response)
			}
		})
	}
}

func Test_handler_register(t *testing.T) {
	type args struct {
		c       *gin.Context
		request []byte
	}
	tests := []struct {
		name               string
		args               args
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "happy path",
			args: args{
				c:       nil,
				request: validRegisterRequestBody,
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: RegisterResponse{
				Success: true,
				Message: "User registered successfully",
			},
		},
		{
			name: "invalid request",
			args: args{
				c:       nil,
				request: invalidRegisterRequestBody,
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Assemble
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			payload := bytes.NewBuffer(tt.args.request)
			r, _ := http.NewRequest("POST", "/auth/register", payload)
			r.Header.Set("Content-Type", "application/json")
			c.Request = r

			h := &Handler{}

			h.register(c)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			if tt.expectedStatusCode == http.StatusOK && tt.expectedResponse != nil {
				var response RegisterResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.expectedResponse, response)
			}
		})
	}
}

func Test_handler_verifyToken(t *testing.T) {
	type args struct {
		c       *gin.Context
		request []byte
	}
	tests := []struct {
		name               string
		args               args
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "happy path",
			args: args{
				c:       nil,
				request: validVerifyTokenRequestBody,
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: VerifyTokenResponse{
				Valid: true,
			},
		},
		{
			name: "invalid request",
			args: args{
				c:       nil,
				request: invalidVerifyTokenRequestBody,
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Assemble
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			payload := bytes.NewBuffer(tt.args.request)
			r, _ := http.NewRequest("POST", "/auth/verify", payload)
			r.Header.Set("Content-Type", "application/json")
			c.Request = r

			h := &Handler{}

			h.verifyToken(c)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			if tt.expectedStatusCode == http.StatusOK && tt.expectedResponse != nil {
				var response VerifyTokenResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.expectedResponse, response)
			}
		})
	}
}

func Test_handler_logout(t *testing.T) {
	tests := []struct {
		name               string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name:               "happy path",
			expectedStatusCode: http.StatusOK,
			expectedResponse: LogoutResponse{
				Success: true,
				Message: "Logged out successfully",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Assemble
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			r, _ := http.NewRequest("POST", "/auth/logout", nil)
			c.Request = r

			h := &Handler{}

			h.logout(c)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			if tt.expectedStatusCode == http.StatusOK && tt.expectedResponse != nil {
				var response LogoutResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.expectedResponse, response)
			}
		})
	}
} 