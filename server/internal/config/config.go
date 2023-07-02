package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	SecretKey     string            `envconfig:"secret_key" required:"true"`
	Debug         bool              `envconfig:"debug" default:"false"`
	MySQLHost     string            `envconfig:"mysql_host" required:"true"`
	MySQLPort     string            `envconfig:"mysql_port" default:"3306"`
	MySQLUsername string            `envconfig:"mysql_username" required:"true"`
	MySQLPassword string            `envconfig:"mysql_password" required:"true"`
	MySQLDatabase string            `envconfig:"mysql_database" required:"true"`
	MySQLOptions  map[string]string `envconfig:"mysql_options" default:"parseTime:true"`
}

func MustParseConfig(prefix ...string) *Config {
	prf := ""
	if len(prefix) > 0 {
		prf = prefix[0]
	}

	cfg, err := ParseConfig(prf)
	if err != nil {
		panic(err)
	}
	return cfg
}

func ParseConfig(prefix string) (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	validate := validator.New()
	err = validate.Struct(cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) DatabaseURI() string {
	var optParams []string
	for opt, val := range c.MySQLOptions {
		optParams = append(optParams, opt+"="+val)
	}
	optStr := strings.Join(optParams, "&")
	uri := fmt.Sprintf("mysql://%s:%s@%s:%s/%s?%s",
		c.MySQLUsername,
		c.MySQLPassword,
		c.MySQLHost,
		c.MySQLPort,
		c.MySQLDatabase,
		optStr)
	return uri
}
