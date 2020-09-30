package bootstrap

type Config struct {
	Port                int    `short:"p" long:"port" description:"port on which app handles http connection" default:"8000"`
	PortsServiceAddress string `long:"ports_service_address" description:"Host and port of ports service" default:"localhost:9000"`
}
