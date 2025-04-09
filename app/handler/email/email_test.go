package email

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"GoMail/app/logic/email"
	"GoMail/app/logic/email/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test request bodies
var (
	validSendEmailRequestBody   = []byte(`{"from":"sender@example.com","to":"recipient@example.com","subject":"Test Subject","body":"Test Body"}`)
	invalidSendEmailRequestBody = []byte(`{"from":"","to":"recipient@example.com","subject":"Test Subject","body":"Test Body"}`)
)

func Test_handler_sendEmail(t *testing.T) {
	type fields struct {
		email *mocks.Email
	}
	type args struct {
		c       *gin.Context
		request []byte
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name:   "happy path",
			fields: fields{email: buildSendEmailMock(true, &email.SendEmailResponse{Success: true}, nil)},
			args: args{
				c:       nil,
				request: validSendEmailRequestBody,
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: &email.SendEmailResponse{
				Success: true,
			},
		},
		{
			name:   "invalid request",
			fields: fields{email: buildSendEmailMock(false, nil, nil)},
			args: args{
				c:       nil,
				request: invalidSendEmailRequestBody,
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:   "error in logic",
			fields: fields{email: buildSendEmailMock(true, nil, errors.New("failed to send email"))},
			args: args{
				c:       nil,
				request: validSendEmailRequestBody,
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Assemble
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			payload := bytes.NewBuffer(tt.args.request)
			r, _ := http.NewRequest("POST", "/email/send", payload)
			r.Header.Set("Content-Type", "application/json")
			c.Request = r

			h := &Handler{
				emailService: tt.fields.email,
			}

			h.sendEmail(c)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			if tt.expectedStatusCode == http.StatusOK && tt.expectedResponse != nil {
				var response email.SendEmailResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.expectedResponse, &response)
			}
			
			tt.fields.email.AssertExpectations(t)
		})
	}
}

func buildSendEmailMock(enableFlag bool, res *email.SendEmailResponse, err error) *mocks.Email {
	client := &mocks.Email{}
	if enableFlag {
		client.On("Send", mock.Anything, mock.AnythingOfType("email.SendEmailRequest")).Return(res, err)
	}
	return client
}

func Test_handler_sendHTMLEmail(t *testing.T) {
	type fields struct {
		email *mocks.Email
	}
	type args struct {
		c       *gin.Context
		request []byte
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name:   "happy path",
			fields: fields{email: buildSendHTMLEmailMock(true, &email.SendEmailResponse{Success: true}, nil)},
			args: args{
				c:       nil,
				request: validSendEmailRequestBody,
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: &email.SendEmailResponse{
				Success: true,
			},
		},
		{
			name:   "invalid request",
			fields: fields{email: buildSendHTMLEmailMock(false, nil, nil)},
			args: args{
				c:       nil,
				request: invalidSendEmailRequestBody,
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:   "error in logic",
			fields: fields{email: buildSendHTMLEmailMock(true, nil, errors.New("failed to send HTML email"))},
			args: args{
				c:       nil,
				request: validSendEmailRequestBody,
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Assemble
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			payload := bytes.NewBuffer(tt.args.request)
			r, _ := http.NewRequest("POST", "/email/send-html", payload)
			r.Header.Set("Content-Type", "application/json")
			c.Request = r

			h := &Handler{
				emailService: tt.fields.email,
			}

			h.sendHTMLEmail(c)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			if tt.expectedStatusCode == http.StatusOK && tt.expectedResponse != nil {
				var response email.SendEmailResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.expectedResponse, &response)
			}
			
			tt.fields.email.AssertExpectations(t)
		})
	}
}

func buildSendHTMLEmailMock(enableFlag bool, res *email.SendEmailResponse, err error) *mocks.Email {
	client := &mocks.Email{}
	if enableFlag {
		client.On("SendHTML", mock.Anything, mock.AnythingOfType("email.SendEmailRequest")).Return(res, err)
	}
	return client
}

func Test_handler_getEmailStatus(t *testing.T) {
	tests := []struct {
		name               string
		expectedStatusCode int
	}{
		{
			name:               "success",
			expectedStatusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Assemble
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			r, _ := http.NewRequest("GET", "/email/status", nil)
			c.Request = r

			h := &Handler{}

			h.getEmailStatus(c)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			var response map[string]string
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, "operational", response["status"])
		})
	}
} 