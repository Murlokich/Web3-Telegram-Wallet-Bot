package config

type Config struct {
	Debug bool `envconfig:"DEBUG" default:"false"`
}
