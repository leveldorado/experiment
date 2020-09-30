package adding

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/leveldorado/experiment/grpc/portspb"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockPortsRepo struct {
	mock.Mock
}

func (m *mockPortsRepo) Save(ctx context.Context, p *portspb.Port) error {
	return m.Called(ctx, p).Error(0)
}

func TestService_Save(t *testing.T) {
	repo := &mockPortsRepo{}
	port := &portspb.Port{
		Id:   uuid.New().String(),
		Name: uuid.New().String(),
	}
	repo.On("Save", mock.Anything, port).Return(nil)
	s := NewService(repo)
	resp, err := s.Save(context.Background(), port)
	require.NoError(t, err)
	require.NotNil(t, resp)
	repo.AssertExpectations(t)
}
