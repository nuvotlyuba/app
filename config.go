package app

type Environment string

type RestServer struct {
	Host string
	Port int
}

type Config struct {
	Env          Environment
	SystemServer RestServer
	ServiceName  string
}

const (
	EnvProd Environment = "prod"
	EnvDev  Environment = "dev"
	EnvTest Environment = "test"
)
