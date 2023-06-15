package app

const (
	appVersion = "v1"
	appName    = "passvault"
	appPort    = ":80"
)

type AppConfig struct {
	appName    string
	appVersion string
	appPort    string
}

func newAppConfig() *AppConfig {
	appConfig := &AppConfig{appName: appName, appVersion: appVersion, appPort: appPort}

	return appConfig
}
