// Package config contains config DAO
package config

// AppConfig is a container for application config
type AppConfig struct {
	ServerConfig *ServerConfig
}

// ServerConfig contains http server settings
type ServerConfig struct {
	Bind string
}

// GetAppConfig returns *Config
func GetAppConfig() *AppConfig {
	bind := ":8080"
	appConf := &AppConfig{
		ServerConfig: &ServerConfig{
			Bind: bind,
		},
	}

	return appConf
}
