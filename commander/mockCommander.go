package commander

import "github.com/stretchr/testify/mock"

// MockCommander ...
type MockCommander struct {
	mock.Mock
}

// GivenMockCommander ...
func GivenMockCommander() *MockCommander {
	return &MockCommander{}
}

// ExecuteCommand ...
func (m *MockCommander) ExecuteCommand(command string, strArgs ...string) (string, error) {
	args := m.Called(command, strArgs)
	return args.String(0), args.Error(1)
}
