package listing

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/leveldorado/experiment/grpc/portspb"
)

type portsRepo interface {
	Get(ctx context.Context, id string) (*portspb.Port, error)
	List(ctx context.Context) (<-chan *portspb.Port, error)
}

/*
Service is responsible for busyness logic
at the moment there is just simple call to repository, but it may be extended with some parameters validation, combining result from several repos e.t.c
*/
type Service struct {
	pr portsRepo
}

func NewService(pr portsRepo) *Service {
	return &Service{pr: pr}
}

func (s *Service) Get(ctx context.Context, req *portspb.GetPortRequest) (*portspb.Port, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid  request")
	}
	port, err := s.pr.Get(ctx, req.Id)
	if err == nil && port == nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf(`port %s is not found`, req.Id))
	}
	return port, err
}

/*
List does not support any parameters like sorting, filtering, but may be extendeds
*/
func (s *Service) List(_ *portspb.ListPortsRequest, stream portspb.ListingService_ListServer) error {
	portsChan, err := s.pr.List(stream.Context())
	if err != nil {
		return err
	}
	for port := range portsChan {
		if err := stream.Send(port); err != nil {
			return err
		}
	}
	return nil
}
