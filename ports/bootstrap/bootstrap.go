package bootstrap

import (
	"fmt"
	"log"
	"net"
	"os"
	"syscall"

	"github.com/jessevdk/go-flags"
	"github.com/leveldorado/experiment/grpc/portspb"
	"github.com/leveldorado/experiment/mongodb"
	"github.com/leveldorado/experiment/ports/adding"
	"github.com/leveldorado/experiment/ports/listing"
	"github.com/leveldorado/experiment/ports/storage"
	"google.golang.org/grpc"
)

type App struct {
	server *grpc.Server
	cfg    Config
}

func (a *App) Shutdown() {
	a.server.GracefulStop()
}

func (a *App) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.cfg.Port))
	if err != nil {
		return fmt.Errorf(`failed to listen on port %d: [err: %w]`, a.cfg.Port, err)
	}
	go func() {
		if err := a.server.Serve(l); err != nil {
			log.Println(fmt.Sprintf(`failed to serve grpc connections: [err: %s]`, err))
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}
	}()
	return nil
}

func (a *App) Build() error {
	cfg, err := getConfig()
	if err != nil {
		return err
	}
	repo, err := getPortsRepo(cfg)
	if err != nil {
		return err
	}
	a.server = getServer(repo)
	a.cfg = cfg
	return nil
}

func getServer(repo *storage.MongodbPortsRepo) *grpc.Server {
	server := grpc.NewServer()
	portspb.RegisterAddingServiceServer(server, adding.NewService(repo))
	portspb.RegisterListingServiceServer(server, listing.NewService(repo))
	return server
}

func getPortsRepo(cfg Config) (*storage.MongodbPortsRepo, error) {
	mongodbClient, err := mongodb.GetClient(cfg.MongodbURL, mongodb.DefaultMongodbConnectTimeout)
	if err != nil {
		return nil, fmt.Errorf(`failed to get mongodb client: [url: %s, err: %w]`, cfg.MongodbURL, err)
	}
	return storage.NewMongodbPortsRepo(mongodbClient, storage.ExperimentDatabaseName, storage.PortCollectionName), nil
}

func getConfig() (Config, error) {
	cfg := Config{}
	/*
	   to simplify config management - config  parameters are taken from command args
	   usually using just command args is not a good option, since we might need a lot of config parameters and have ability to change them without app restart
	*/
	if _, err := flags.ParseArgs(&cfg, os.Args); err != nil {
		return Config{}, fmt.Errorf(`failed to parse args: [args: %v, err: %w]`, os.Args, err)
	}
	return cfg, nil
}
