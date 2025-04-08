// email/send_html_test.go
//go:build unit
// +build unit

package email

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"GoMail/app/config"
	"GoMail/app/libs/smtp/mocks"
	repoMocks "GoMail/app/repository/mocks"
)

var (
	validSendHTMLEmailRequest = SendEmailRequest{
		From:    "test@example.com",
		To:      "recipient@example.com",
		Subject: "Test Subject",
		Body:    "<h1>Test HTML Body</h1>",
	}
	
	htmlSuccessResponse = &SendEmailResponse{
		Success: true,
		Error:   "",
	}
	
	htmlErrorResponse = &SendEmailResponse{
		Success: false,
		Error:   "smtp error",
	}
	
	htmlTestError = errors.New("smtp error")
)

func buildHTMLMockSMTPClient(success bool) *mocks.SMTPClient {
	client := &mocks.SMTPClient{}
	if success {
		client.On("SendHTML", mock.Anything, "test@example.com", "recipient@example.com", "Test Subject", "<h1>Test HTML Body</h1>").Return(nil)
	} else {
		client.On("SendHTML", mock.Anything, "test@example.com", "recipient@example.com", "Test Subject", "<h1>Test HTML Body</h1>").Return(htmlTestError)
	}
	return client
}

func buildHTMLMockRepo() *repoMocks.Repository {
	repo := &repoMocks.Repository{}
	repo.On("SaveEmailLog", mock.Anything, mock.AnythingOfType("*models.EmailLog")).Return(nil)
	return repo
}

func TestEmailService_SendHTML(t *testing.T) {
	type fields struct {
		client *mocks.SMTPClient
		repo   *repoMocks.Repository
		config *config.Config
	}
	
	type args struct {
		ctx context.Context
		req SendEmailRequest
	}
	
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *SendEmailResponse
		wantErr error
	}{
		{
			name: "happy path",
			fields: fields{
				client: buildHTMLMockSMTPClient(true),
				repo:   buildHTMLMockRepo(),
				config: &config.Config{},
			},
			args: args{
				ctx: context.Background(),
				req: validSendHTMLEmailRequest,
			},
			want:    htmlSuccessResponse,
			wantErr: nil,
		},
		{
			name: "smtp error",
			fields: fields{
				client: buildHTMLMockSMTPClient(false),
				repo:   buildHTMLMockRepo(),
				config: &config.Config{},
			},
			args: args{
				ctx: context.Background(),
				req: validSendHTMLEmailRequest,
			},
			want:    htmlErrorResponse,
			wantErr: htmlTestError,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &emailService{
				client: tt.fields.client,
				repo:   tt.fields.repo,
				config: tt.fields.config,
			}
			
			got, err := s.SendHTML(tt.args.ctx, tt.args.req)
			
			// Add a small delay to allow the goroutine to complete
			time.Sleep(100 * time.Millisecond)
			
			assert.Equal(t, tt.want, got)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error()) 
			} else {
				assert.NoError(t, err)
			}
			
			tt.fields.client.AssertExpectations(t)
			tt.fields.repo.AssertExpectations(t)
		})
	}
} 