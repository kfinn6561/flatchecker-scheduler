package testutil

import (
	"context"
	"flatchecker-scheduler/pubsublib"

	"github.com/stretchr/testify/mock"
)

// MockSecretGetter is a mock implementation of the SecretGetter interface
type MockSecretGetter struct {
	mock.Mock
}

func (m *MockSecretGetter) GetSecret(name string) (string, error) {
	args := m.Called(name)
	return args.String(0), args.Error(1)
}

// MockPublisher is a mock implementation of the Publisher interface
type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) PublishSchedules(ctx context.Context, schedules []pubsublib.ScheduledSearchesMessage) error {
	args := m.Called(ctx, schedules)
	return args.Error(0)
}
