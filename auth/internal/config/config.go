package config

import "github.com/singl3focus/cli_messenger/auth/internal/config/env"

type Config interface {
	Load(path string) error

	GRPCConfig
	PGConfig
}

type GRPCConfig interface {
	GRPCPort() int
}

type PGConfig interface {
	PGDSN() string 
}

type ConfigType int

const (
	ENV ConfigType = iota
	YAML
	JSON
)

func NewConfig(option ConfigType) Config {
	switch option {
	case ENV:
		return env.NewConfig()
	case YAML:
		fallthrough
	case JSON:
		panic("not implemented")
	}
	return nil
}