// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/kujilabo/cocotola-synthesizer-api/src/app/domain"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// SynthesizerClient is an autogenerated mock type for the SynthesizerClient type
type SynthesizerClient struct {
	mock.Mock
}

// Synthesize provides a mock function with given fields: ctx, lang5, text
func (_m *SynthesizerClient) Synthesize(ctx context.Context, lang5 domain.Lang5, text string) (string, error) {
	ret := _m.Called(ctx, lang5, text)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, domain.Lang5, string) string); ok {
		r0 = rf(ctx, lang5, text)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.Lang5, string) error); ok {
		r1 = rf(ctx, lang5, text)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSynthesizerClient creates a new instance of SynthesizerClient. It also registers a cleanup function to assert the mocks expectations.
func NewSynthesizerClient(t testing.TB) *SynthesizerClient {
	mock := &SynthesizerClient{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}