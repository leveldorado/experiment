package adding

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/leveldorado/experiment/grpc/portspb"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type mockAddingServiceClient struct {
	mock.Mock
}

func (m *mockAddingServiceClient) Save(ctx context.Context, in *portspb.Port, opts ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*empty.Empty), args.Error(1)
}

func TestMakePOSTPortEndpoint(t *testing.T) {
	js, cl := generatePortsJSONAndPreparePortsServiceClient(t)
	r := httprouter.New()
	RegisterEndpoints(r, cl)
	req, err := http.NewRequest(http.MethodPost, "/api/v1/ports", js)
	require.NoError(t, err)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)
	cl.AssertExpectations(t)
}

func generatePortsJSONAndPreparePortsServiceClient(t *testing.T) (io.Reader, *mockAddingServiceClient) {
	cl := &mockAddingServiceClient{}
	portsMap := map[string]*portspb.Port{}
	for i := 0; i < 100; i++ {
		port := generatePort()
		portsMap[port.Id] = port
		cl.On("Save", mock.Anything, port, mock.Anything).Return(&empty.Empty{}, nil).Once()
	}
	buff := &bytes.Buffer{}
	err := json.NewEncoder(buff).Encode(portsMap)
	require.NoError(t, err)
	return buff, cl
}

func generatePort() *portspb.Port {
	return &portspb.Port{
		Id:          uuid.New().String(),
		Name:        uuid.New().String(),
		City:        uuid.New().String(),
		Country:     uuid.New().String(),
		Alias:       []string{uuid.New().String()},
		Regions:     []string{uuid.New().String()},
		Coordinates: []float32{12, 45.01},
		Province:    uuid.New().String(),
		Timezone:    uuid.New().String(),
		Unlocks:     []string{uuid.New().String()},
	}
}
