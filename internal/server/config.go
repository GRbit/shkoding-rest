package server

// nolint:gochecknoglobals
var (
	Version string
	Built   string
)

type serviceConfig struct {
	Service        instanceConfig
}

type instanceConfig struct {
	Addr        string `long:"addr" default:"localhost:1444" env:"SERVICE_ADDR" description:"service address"`
	SystemToken string `long:"token" default:"123" env:"SYSTEM_TOKEN" description:"token for clients authorization"`

	Debug       bool   `long:"debug" env:"SERVICE_DEBUG" description:"sets debug mode logging, overrides log-level"`
	Console     bool   `long:"console" env:"CONSOLE" description:"extended debug mode; adapts logs output for console"`
	LogLevel    string `long:"log-level" default:"info" env:"LOG_LEVEL" description:"set log level (debug|info|warn |error)"`
	Name        string `long:"service-name" default:"vk-stats" env:"SERVICE_NAME" description:"service name for logging"`
}
