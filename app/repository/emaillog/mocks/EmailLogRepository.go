// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"
	models "GoMail/app/repository/models"
	mock "github.com/stretchr/testify/mock"
)

// EmailLogRepository is an autogenerated mock type for the EmailLogRepository type
type EmailLogRepository struct {
	mock.Mock
}

// FindAll provides a mock function with given fields: ctx, filter, page, limit
func (_m *EmailLogRepository) FindAll(ctx context.Context, filter interface{}, page int, limit int) ([]*models.EmailLog, int64, error) {
	ret := _m.Called(ctx, filter, page, limit)

	var r0 []*models.EmailLog
	var r1 int64
	var r2 error

	if rf, ok := ret.Get(0).(func(context.Context, interface{}, int, int) ([]*models.EmailLog, int64, error)); ok {
		return rf(ctx, filter, page, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, int, int) []*models.EmailLog); ok {
		r0 = rf(ctx, filter, page, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.EmailLog)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, interface{}, int, int) int64); ok {
		r1 = rf(ctx, filter, page, limit)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(context.Context, interface{}, int, int) error); ok {
		r2 = rf(ctx, filter, page, limit)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// FindByID provides a mock function with given fields: ctx, id
func (_m *EmailLogRepository) FindByID(ctx context.Context, id string) (*models.EmailLog, error) {
	ret := _m.Called(ctx, id)

	var r0 *models.EmailLog
	var r1 error

	if rf, ok := ret.Get(0).(func(context.Context, string) (*models.EmailLog, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.EmailLog); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.EmailLog)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, emailLog
func (_m *EmailLogRepository) Save(ctx context.Context, emailLog *models.EmailLog) error {
	ret := _m.Called(ctx, emailLog)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.EmailLog) error); ok {
		r0 = rf(ctx, emailLog)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewEmailLogRepository creates a new instance of EmailLogRepository
func NewEmailLogRepository() *EmailLogRepository {
	return &EmailLogRepository{}
} 