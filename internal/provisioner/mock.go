package provisioner

import (
	"context"
	"github.com/broswen/vex/internal/flag"
	"github.com/broswen/vex/internal/project"
	"github.com/broswen/vex/internal/token"
	"github.com/stretchr/testify/mock"
)

type MockProvisioner struct {
	mock.Mock
}

func NewMockProvisioner() *MockProvisioner {
	return &MockProvisioner{}
}

func (m *MockProvisioner) ProvisionFlag(ctx context.Context, f *flag.Flag) error {
	args := m.Called(ctx, f)
	return args.Error(0)
}

func (m *MockProvisioner) DeprovisionFlag(ctx context.Context, f *flag.Flag) error {
	args := m.Called(ctx, f)
	return args.Error(0)
}

func (m *MockProvisioner) ProvisionProject(ctx context.Context, p *project.Project) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockProvisioner) DeprovisionProject(ctx context.Context, p *project.Project) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockProvisioner) ProvisionToken(ctx context.Context, t *token.Token) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *MockProvisioner) DeprovisionToken(ctx context.Context, t *token.Token) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}
