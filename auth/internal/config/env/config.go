package env

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	grpcPort = "GRPC_PORT"

	postgresDSN = "POSTGRES_DSN"
)

type Config struct {}

func NewConfig() Config {
	return Config{}
}

func (c Config) Load(path string) error {
	return godotenv.Load(path)
}

func (c Config) GRPCPort() int {	
	portString := os.Getenv(grpcPort)
	if portString == "" {
		panic("var is empty: " + grpcPort)
	}

	portInt, err := strconv.Atoi(grpcPort)
	if err != nil {
		panic(err)
	}

	return portInt
}

func (c Config) PGDSN() string {
	dsn := os.Getenv(postgresDSN)
	if dsn == "" {
		panic("var is empty: " + postgresDSN)
	}

	return dsn
}
