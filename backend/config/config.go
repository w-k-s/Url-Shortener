package config

import (
	env "github.com/Netflix/go-env"
	"github.com/w-k-s/short-url/log"
	"net/url"
)

type settings struct {
	DatabaseConnectionString       string `env:"DB_CONN_STRING,required=true"`
	ListenAddress                  string `env:"ADDRESS,default=:80"`
	BaseURL                        string `env:"BASE_URL,required=true"`
	AccessControlAllowOriginHeader string `env:"ALLOW_ORIGIN"`
	baseURL                        *url.URL
}

var Settings settings
var initialized bool

func (s settings) GetBaseURL() *url.URL {
	return s.baseURL
}

func Init() {
	_, err := env.UnmarshalFromEnviron(&Settings)
	if err != nil {
		log.Fatal(err)
	}

	baseURL, err := url.Parse(Settings.BaseURL)
	if err != nil {
		log.Fatalf("Failed to parse env variable 'BASE_URL': '%s'", Settings.BaseURL)
	}
	if len(baseURL.Scheme) == 0 {
		log.Fatalf("Failed to determine scheme from BASE_URL %q", baseURL)
	}
	if len(baseURL.Host) == 0 {
		log.Fatalf("Failed to determine host from BASE_URL %q", baseURL)
	}
	Settings.baseURL = baseURL
}
