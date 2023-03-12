package app

const (
	appVersion = "v1"
	appName    = "passvault"
	appPort    = ":80"
)

type AppConfig struct {
	AppName    string
	AppVersion string
	AppPort    string
}

func GetAppConfig() *AppConfig {
	appConfig := &AppConfig{AppName: appName, AppVersion: appVersion, AppPort: appPort}

	return appConfig
}
