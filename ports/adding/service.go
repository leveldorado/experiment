package adding

import (
	"context"
	"fmt"

	"github.com/leveldorado/experiment/grpc/portspb"

	"github.com/golang/protobuf/ptypes/empty"
)

type portsRepo interface {
	Save(ctx context.Context, p *portspb.Port) error
}

/*
Service is responsible for busyness logic
at the moment there is just simple call to repository, but it may be extended with some validation, firing events, audits e.t.c
*/
type Service struct {
	portspb.UnimplementedAddingServiceServer
	pr portsRepo
}

func NewService(pr portsRepo) *Service {
	return &Service{pr: pr}
}

func (s *Service) Save(ctx context.Context, p *portspb.Port) (*empty.Empty, error) {
	if err := s.pr.Save(ctx, p); err != nil {
		return nil, fmt.Errorf(`failed to save port: [port: %+v, err: %w]`, p, err)
	}
	return &empty.Empty{}, nil
}
