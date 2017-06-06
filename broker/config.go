package broker

type Config struct {
	RedisURL      string
	RedisPassword string
}

func NewConfig() *Config {

	return &Config{}
}
