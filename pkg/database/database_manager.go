package database

import (
	"database/sql"
	"fmt"
	"log"
	"passvault/pkg/types"

	_ "github.com/go-sql-driver/mysql"
)

var (
	databaseManager *DatabaseManager
)

type DatabaseManager struct {
	dbConnection *sql.DB
}

func formatCredentials(databaseConfig DatabaseConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", databaseConfig.Username, databaseConfig.Password, databaseConfig.Host, databaseConfig.Port, databaseConfig.DatabaseName)
}

func Get() *DatabaseManager {
	if databaseManager == nil {
		log.Println("creds are", formatCredentials(*GetDatabaseConfig()))
		dbConn, err := sql.Open("mysql", formatCredentials(*GetDatabaseConfig()))
		if err != nil {
			panic(err)
		}

		databaseManager = &DatabaseManager{dbConnection: dbConn}
	}

	return databaseManager
}

func (dm *DatabaseManager) Save(entry types.Entry) error {
	stmt, err := dm.dbConnection.Prepare("insert into db.passwords (domain, username, password) VALUES (?, ?, ?)")
	if err != nil {
		log.Println(err)
		return err
	}

	defer stmt.Close()

	log.Println("entryto e", entry)

	_, err = stmt.Exec(entry.Domain, entry.Username, entry.Password)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (dm *DatabaseManager) Get(domain string) {

}

func (dm *DatabaseManager) Contains(domain string) {

}
