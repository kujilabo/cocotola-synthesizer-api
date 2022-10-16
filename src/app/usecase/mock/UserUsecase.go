// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/kujilabo/cocotola-synthesizer-api/src/app/domain"
	mock "github.com/stretchr/testify/mock"

	service "github.com/kujilabo/cocotola-synthesizer-api/src/app/service"

	testing "testing"
)

// UserUsecase is an autogenerated mock type for the UserUsecase type
type UserUsecase struct {
	mock.Mock
}

// FindAudioByAudioID provides a mock function with given fields: ctx, audioID
func (_m *UserUsecase) FindAudioByAudioID(ctx context.Context, audioID domain.AudioID) (service.Audio, error) {
	ret := _m.Called(ctx, audioID)

	var r0 service.Audio
	if rf, ok := ret.Get(0).(func(context.Context, domain.AudioID) service.Audio); ok {
		r0 = rf(ctx, audioID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.Audio)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.AudioID) error); ok {
		r1 = rf(ctx, audioID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Synthesize provides a mock function with given fields: ctx, lang2, text
func (_m *UserUsecase) Synthesize(ctx context.Context, lang2 domain.Lang2, text string) (service.Audio, error) {
	ret := _m.Called(ctx, lang2, text)

	var r0 service.Audio
	if rf, ok := ret.Get(0).(func(context.Context, domain.Lang2, string) service.Audio); ok {
		r0 = rf(ctx, lang2, text)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service.Audio)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.Lang2, string) error); ok {
		r1 = rf(ctx, lang2, text)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUserUsecase creates a new instance of UserUsecase. It also registers a cleanup function to assert the mocks expectations.
func NewUserUsecase(t testing.TB) *UserUsecase {
	mock := &UserUsecase{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
