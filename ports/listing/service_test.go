package listing

import (
	"context"
	"testing"

	"google.golang.org/grpc"

	"github.com/google/uuid"
	"github.com/leveldorado/experiment/grpc/portspb"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockPortsRepo struct {
	mock.Mock
}

func (m *mockPortsRepo) Get(ctx context.Context, id string) (*portspb.Port, error) {
	args := m.Called(ctx, id)
	port, _ := args.Get(0).(*portspb.Port)
	return port, args.Error(1)
}
func (m *mockPortsRepo) List(ctx context.Context) (<-chan *portspb.Port, error) {
	args := m.Called(ctx)
	ch, _ := args.Get(0).(<-chan *portspb.Port)
	return ch, args.Error(1)
}

func TestService_Get(t *testing.T) {
	repo := &mockPortsRepo{}
	port := &portspb.Port{
		Id:   uuid.New().String(),
		Name: uuid.New().String(),
	}
	notFoundID := "not found"
	repo.On("Get", mock.Anything, port.Id).Return(port, nil)
	repo.On("Get", mock.Anything, notFoundID).Return(nil, nil)
	s := NewService(repo)
	respPort, err := s.Get(context.Background(), &portspb.GetPortRequest{Id: port.Id})
	require.NoError(t, err)
	require.Equal(t, port, respPort)

	_, err = s.Get(context.Background(), &portspb.GetPortRequest{Id: notFoundID})
	require.Equal(t, codes.NotFound, status.Code(err))
	repo.AssertExpectations(t)
}

type mockListingService_ListServer struct {
	grpc.ServerStream
	mock.Mock
}

func (m *mockListingService_ListServer) Send(p *portspb.Port) error {
	return m.Called(p).Error(0)
}

func (m *mockListingService_ListServer) Context() context.Context {
	return m.Called().Get(0).(context.Context)
}

func TestNewService(t *testing.T) {
	stream := &mockListingService_ListServer{}
	ports := []*portspb.Port{{Id: uuid.New().String()}, {Id: uuid.New().String()}}
	for _, p := range ports {
		stream.On("Send", p).Return(nil).Once()
	}
	stream.On("Context").Return(context.Background())
	repo := &mockPortsRepo{}
	repo.On("List", mock.Anything).Return(getPortsChan(ports), nil)
	s := NewService(repo)
	require.NoError(t, s.List(nil, stream))
	repo.AssertExpectations(t)
	stream.AssertExpectations(t)
}

func getPortsChan(ports []*portspb.Port) <-chan *portspb.Port {
	portsChan := make(chan *portspb.Port)
	go func() {
		for _, port := range ports {
			portsChan <- port
		}
		close(portsChan)
	}()
	return portsChan
}
