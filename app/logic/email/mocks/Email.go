// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"
	email "GoMail/app/logic/email"
	mock "github.com/stretchr/testify/mock"
)

// Email is an autogenerated mock type for the Email type
type Email struct {
	mock.Mock
}

// Send provides a mock function with given fields: ctx, req
func (_m *Email) Send(ctx context.Context, req email.SendEmailRequest) (*email.SendEmailResponse, error) {
	ret := _m.Called(ctx, req)

	var r0 *email.SendEmailResponse
	var r1 error

	if rf, ok := ret.Get(0).(func(context.Context, email.SendEmailRequest) (*email.SendEmailResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, email.SendEmailRequest) *email.SendEmailResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*email.SendEmailResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, email.SendEmailRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendBulk provides a mock function with given fields: ctx, req
func (_m *Email) SendBulk(ctx context.Context, req email.SendBulkEmailRequest) (*email.SendBulkEmailResponse, error) {
	ret := _m.Called(ctx, req)

	var r0 *email.SendBulkEmailResponse
	var r1 error

	if rf, ok := ret.Get(0).(func(context.Context, email.SendBulkEmailRequest) (*email.SendBulkEmailResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, email.SendBulkEmailRequest) *email.SendBulkEmailResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*email.SendBulkEmailResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, email.SendBulkEmailRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendHTML provides a mock function with given fields: ctx, req
func (_m *Email) SendHTML(ctx context.Context, req email.SendEmailRequest) (*email.SendEmailResponse, error) {
	ret := _m.Called(ctx, req)

	var r0 *email.SendEmailResponse
	var r1 error

	if rf, ok := ret.Get(0).(func(context.Context, email.SendEmailRequest) (*email.SendEmailResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, email.SendEmailRequest) *email.SendEmailResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*email.SendEmailResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, email.SendEmailRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendWithAttachments provides a mock function with given fields: ctx, req
func (_m *Email) SendWithAttachments(ctx context.Context, req email.SendWithAttachmentsRequest) (*email.SendEmailResponse, error) {
	ret := _m.Called(ctx, req)

	var r0 *email.SendEmailResponse
	var r1 error

	if rf, ok := ret.Get(0).(func(context.Context, email.SendWithAttachmentsRequest) (*email.SendEmailResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, email.SendWithAttachmentsRequest) *email.SendEmailResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*email.SendEmailResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, email.SendWithAttachmentsRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewEmail creates a new instance of Email
func NewEmail() *Email {
	return &Email{}
} 