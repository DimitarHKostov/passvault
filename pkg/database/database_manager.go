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
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(entry.Domain, entry.Username, entry.Password)
	if err != nil {
		return err
	}

	return nil
}

func (dm *DatabaseManager) Get(domain string) (*types.Entry, error) {
	stmt, err := dm.dbConnection.Prepare("SELECT * FROM passwords WHERE domain = ?")
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	var entry types.Entry

	row := stmt.QueryRow(domain)

	err = row.Scan(&entry.Domain, &entry.Username, &entry.Password)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (dm *DatabaseManager) Contains(domain string) (bool, error) {
	query := "SELECT COUNT(*) FROM passwords WHERE domain = ?"

	stmt, err := dm.dbConnection.Prepare(query)
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	var count int
	err = stmt.QueryRow(domain).Scan(&count)
	if err != nil {
		return false, err
	}

	return count != 0, nil
}
