package bootstrap

type Config struct {
	MongodbURL string `short:"m" long:"mongodb_url" description:"Mongodb connect url" default:"mongodb://localhost:27017"`
	Port       int    `short:"p" long:"port" description:"port on which app handles grpc connection" default:"9000"`
}
