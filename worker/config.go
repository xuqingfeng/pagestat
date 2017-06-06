package worker

type Config struct {
	RedisUrl      string
	RedisPassword string
}

func NewConfig() *Config {

	return &Config{}
}
