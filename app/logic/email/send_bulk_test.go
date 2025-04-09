package email

import (
	"context"
	"errors"
	"testing"
	"time"

	"GoMail/app/config"
	"GoMail/app/libs/smtp/mocks"
	repoMocks "GoMail/app/repository/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	validBulkEmail1 = BulkEmail{
		From:    "test@example.com",
		To:      "recipient1@example.com",
		Subject: "Test Subject 1",
		Body:    "Test Body 1",
		IsHTML:  false,
	}
	
	validBulkEmail2 = BulkEmail{
		From:    "test@example.com",
		To:      "recipient2@example.com",
		Subject: "Test Subject 2",
		Body:    "<h1>Test Body 2</h1>",
		IsHTML:  true,
	}
	
	validBulkEmailRequest = SendBulkEmailRequest{
		Emails: []BulkEmail{validBulkEmail1, validBulkEmail2},
	}
	
	bulkSuccessResponse = &SendBulkEmailResponse{
		Results: []EmailResult{
			{Success: true, Error: ""},
			{Success: true, Error: ""},
		},
	}
	
	bulkPartialFailureResponse = &SendBulkEmailResponse{
		Results: []EmailResult{
			{Success: true, Error: ""},
			{Success: false, Error: "smtp error"},
		},
	}
	
	bulkAllFailuresResponse = &SendBulkEmailResponse{
		Results: []EmailResult{
			{Success: false, Error: "smtp error 1"},
			{Success: false, Error: "smtp error 2"},
		},
	}
	
	testError1 = errors.New("smtp error 1")
	testError2 = errors.New("smtp error 2")
)

func buildBulkMockSMTPClient(status string) *mocks.SMTPClient {
	client := &mocks.SMTPClient{}
	
	switch status {
	case "success":
		client.On("Send", mock.Anything, "test@example.com", "recipient1@example.com", "Test Subject 1", "Test Body 1").Return(nil)
		client.On("SendHTML", mock.Anything, "test@example.com", "recipient2@example.com", "Test Subject 2", "<h1>Test Body 2</h1>").Return(nil)
	case "partial_failure":
		client.On("Send", mock.Anything, "test@example.com", "recipient1@example.com", "Test Subject 1", "Test Body 1").Return(nil)
		client.On("SendHTML", mock.Anything, "test@example.com", "recipient2@example.com", "Test Subject 2", "<h1>Test Body 2</h1>").Return(errors.New("smtp error"))
	case "all_failures":
		client.On("Send", mock.Anything, "test@example.com", "recipient1@example.com", "Test Subject 1", "Test Body 1").Return(testError1)
		client.On("SendHTML", mock.Anything, "test@example.com", "recipient2@example.com", "Test Subject 2", "<h1>Test Body 2</h1>").Return(testError2)
	}
	
	return client
}

func buildBulkMockRepo() *repoMocks.Repository {
	repo := &repoMocks.Repository{}
	repo.On("SaveEmailLog", mock.Anything, mock.AnythingOfType("*models.EmailLog")).Return(nil)
	return repo
}

func TestEmailService_SendBulk(t *testing.T) {
	type fields struct {
		client *mocks.SMTPClient
		repo   *repoMocks.Repository
		config *config.Config
	}
	
	type args struct {
		ctx context.Context
		req SendBulkEmailRequest
	}
	
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *SendBulkEmailResponse
		wantErr error
	}{
		{
			name: "happy path",
			fields: fields{
				client: buildBulkMockSMTPClient("success"),
				repo:   buildBulkMockRepo(),
				config: &config.Config{},
			},
			args: args{
				ctx: context.Background(),
				req: validBulkEmailRequest,
			},
			want:    bulkSuccessResponse,
			wantErr: nil,
		},
		{
			name: "partial failure",
			fields: fields{
				client: buildBulkMockSMTPClient("partial_failure"),
				repo:   buildBulkMockRepo(),
				config: &config.Config{},
			},
			args: args{
				ctx: context.Background(),
				req: validBulkEmailRequest,
			},
			want:    bulkPartialFailureResponse,
			wantErr: nil,
		},
		{
			name: "all failures",
			fields: fields{
				client: buildBulkMockSMTPClient("all_failures"),
				repo:   buildBulkMockRepo(),
				config: &config.Config{},
			},
			args: args{
				ctx: context.Background(),
				req: validBulkEmailRequest,
			},
			want:    bulkAllFailuresResponse,
			wantErr: nil,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &emailService{
				client: tt.fields.client,
				repo:   tt.fields.repo,
				config: tt.fields.config,
			}
			
			got, err := s.SendBulk(tt.args.ctx, tt.args.req)
			
			// Add a small delay to allow the goroutines to complete
			time.Sleep(200 * time.Millisecond)
			
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			
			// Verify results array length matches
			assert.Equal(t, len(tt.want.Results), len(got.Results))
			
			// Verify success field for each result
			for i, wantResult := range tt.want.Results {
				assert.Equal(t, wantResult.Success, got.Results[i].Success)
				
				// For failures, verify error message
				if !wantResult.Success {
					assert.Equal(t, wantResult.Error, got.Results[i].Error)
				}
			}
			
			// Verify mock expectations
			tt.fields.client.AssertExpectations(t)
			tt.fields.repo.AssertExpectations(t)
		})
	}
} 