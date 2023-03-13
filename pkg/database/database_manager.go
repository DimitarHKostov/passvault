package database

import (
	"database/sql"
	"fmt"
	"log"
	"passvault/pkg/types"

	_ "github.com/lib/pq"
)

var (
	databaseManager *DatabaseManager
)

type DatabaseManager struct {
	dbConnection *sql.DB
}

func formatCredentials(databaseConfig DatabaseConfig) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", databaseConfig.Host, databaseConfig.Port, databaseConfig.Username, databaseConfig.Password, databaseConfig.DatabaseName)
}

func Get() *DatabaseManager {
	if databaseManager == nil {
		dbConnection, err := sql.Open("postgres", formatCredentials(*GetDatabaseConfig()))
		if err != nil {
			panic(err)
		}

		databaseManager = &DatabaseManager{dbConnection: dbConnection}
	}

	return databaseManager
}

func (dm *DatabaseManager) Save(entry types.Entry) error {
	query := `insert into passwords(domain, username, password) VALUES ($1, $2, $3)`

	_, err := dm.dbConnection.Exec(query, entry.Domain, entry.Username, entry.Password)
	if err != nil {
		log.Println("tuka?")
		return err
	}

	return nil
}

func (dm *DatabaseManager) Get(domain string) {

}

func (dm *DatabaseManager) Contains(domain string) {

}
