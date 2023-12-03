package config

type Config struct {
	DBURL      string
	TwitterCFG TwitterConfig
}

type TwitterConfig struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	TokenSecret    string
	BearerToken    string
	BaseURL        string
}

func NewConfig() *Config {
	return &Config{}
}
