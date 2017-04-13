package worker

type Config struct {
	NsqLookupdAddr string
}

func NewConfig() *Config {

	return &Config{}
}
