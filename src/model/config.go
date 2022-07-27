package model

type Config struct {
	Services []Service `mapstructure:"services"`
	Auth     Auth      `mapstructure:"auth"`
}

type Service struct {
	Name       string            `mapstructure:"name"`
	Prefix     string            `mapstructure:"prefix"`
	ServiceUrl string            `mapstructure:"serviceUrl"`
	Headers    map[string]string `mapstructure:"headers"`
	PublicUrls []string          `mapstructure:"publicUrls"`
}

type Auth struct {
	LoginUrl         string `mapstructure:"loginUrl"`
	ValidateTokenUrl string `mapstructure:"validateTokenUrl"`
	UserHubApiKey    string `mapstructure:"userHubApiKey"`
}
