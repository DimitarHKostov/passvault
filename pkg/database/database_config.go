package database

const (
	host         = "localhost"
	port         = "5432"
	username     = "postgres"
	password     = "mysecretpassword"
	databaseName = "db"
)

var (
	databaseConfig *DatabaseConfig
)

type DatabaseConfig struct {
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
}

func GetDatabaseConfig() *DatabaseConfig {
	if databaseConfig == nil {
		databaseConfig = &DatabaseConfig{
			Host:         host,
			Port:         port,
			Username:     username,
			Password:     password,
			DatabaseName: databaseName,
		}
	}

	return databaseConfig
}
