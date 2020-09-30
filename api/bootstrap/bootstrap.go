package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"

	"github.com/jessevdk/go-flags"
	"github.com/julienschmidt/httprouter"
	"github.com/leveldorado/experiment/api/adding"
	"github.com/leveldorado/experiment/api/listing"
	"github.com/leveldorado/experiment/grpc/portspb"
	"google.golang.org/grpc"
)

type App struct {
	server   *http.Server
	grpcConn *grpc.ClientConn
}

func (a *App) Shutdown() error {
	if err := a.grpcConn.Close(); err != nil {
		return err
	}
	return a.server.Shutdown(context.Background())
}

func (a *App) Run() error {
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println(fmt.Sprintf(`failed to serve http connections: [err: %s]`, err))
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
	conn, err := grpc.Dial(cfg.PortsServiceAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}
	router := httprouter.New()
	adding.RegisterEndpoints(router, portspb.NewAddingServiceClient(conn))
	listing.RegisterEndpoints(router, portspb.NewListingServiceClient(conn))
	a.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}
	a.grpcConn = conn
	return nil
}

func getConfig() (Config, error) {
	cfg := Config{}
	/*
	   to simplify config management - config  parameters are taken from command args
	   usually using just command args is not a good option, since we might need a lot of config parameters and have ability to change them without app restart
	*/
	if _, err := flags.ParseArgs(&cfg, os.Args[1:]); err != nil {
		return Config{}, fmt.Errorf(`failed to parse args: [args: %v, err: %w]`, os.Args, err)
	}
	return cfg, nil
}
