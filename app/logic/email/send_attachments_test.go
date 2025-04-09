package email

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"GoMail/app/config"
	libSmtp "GoMail/app/libs/smtp"
	"GoMail/app/libs/smtp/mocks"
	repoMocks "GoMail/app/repository/mocks"
)

var (
	validAttachment = libSmtp.Attachment{
		Filename: "test.txt",
		Content:  []byte("test content"),
		MimeType: "text/plain",
	}
	
	validSendWithAttachmentsRequest = SendWithAttachmentsRequest{
		From:        "test@example.com",
		To:          "recipient@example.com",
		Subject:     "Test Subject",
		Body:        "Test Body",
		Attachments: []libSmtp.Attachment{validAttachment},
	}
	
	attachmentSuccessResponse = &SendEmailResponse{
		Success: true,
		Error:   "",
	}
	
	attachmentErrorResponse = &SendEmailResponse{
		Success: false,
		Error:   "smtp error",
	}
	
	attachmentTestError = errors.New("smtp error")
)

func buildAttachmentsMockSMTPClient(success bool) *mocks.SMTPClient {
	client := &mocks.SMTPClient{}
	if success {
		client.On("SendWithAttachments", 
			mock.Anything, 
			"test@example.com", 
			"recipient@example.com", 
			"Test Subject", 
			"Test Body", 
			[]libSmtp.Attachment{validAttachment}).Return(nil)
	} else {
		client.On("SendWithAttachments", 
			mock.Anything, 
			"test@example.com", 
			"recipient@example.com", 
			"Test Subject", 
			"Test Body", 
			[]libSmtp.Attachment{validAttachment}).Return(attachmentTestError)
	}
	return client
}

func buildAttachmentsMockRepo() *repoMocks.Repository {
	repo := &repoMocks.Repository{}
	repo.On("SaveEmailLog", mock.Anything, mock.AnythingOfType("*models.EmailLog")).Return(nil)
	return repo
}

func TestEmailService_SendWithAttachments(t *testing.T) {
	type fields struct {
		client *mocks.SMTPClient
		repo   *repoMocks.Repository
		config *config.Config
	}
	
	type args struct {
		ctx context.Context
		req SendWithAttachmentsRequest
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
				client: buildAttachmentsMockSMTPClient(true),
				repo:   buildAttachmentsMockRepo(),
				config: &config.Config{},
			},
			args: args{
				ctx: context.Background(),
				req: validSendWithAttachmentsRequest,
			},
			want:    attachmentSuccessResponse,
			wantErr: nil,
		},
		{
			name: "smtp error",
			fields: fields{
				client: buildAttachmentsMockSMTPClient(false),
				repo:   buildAttachmentsMockRepo(),
				config: &config.Config{},
			},
			args: args{
				ctx: context.Background(),
				req: validSendWithAttachmentsRequest,
			},
			want:    attachmentErrorResponse,
			wantErr: attachmentTestError,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &emailService{
				client: tt.fields.client,
				repo:   tt.fields.repo,
				config: tt.fields.config,
			}
			
			got, err := s.SendWithAttachments(tt.args.ctx, tt.args.req)
			
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