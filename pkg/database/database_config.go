package database

type DatabaseConfig struct {
	host         string
	port         string
	username     string
	password     string
	databaseName string
}

func NewDatabaseConfig(host, port, username, password, databaseName string) *DatabaseConfig {
	databaseConfig := &DatabaseConfig{
		host:         host,
		port:         port,
		username:     username,
		password:     password,
		databaseName: databaseName,
	}

	return databaseConfig
}
