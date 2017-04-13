package broker

type Config struct {
	NsqdAddr string
}

func NewConfig() *Config {

	return &Config{}
}
