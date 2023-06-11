package database

import (
	"database/sql"
	"fmt"

	"passvault/pkg/log"
	"passvault/pkg/types"

	_ "github.com/go-sql-driver/mysql"
)

const (
	mysqlDriverName             = "mysql"
	mysqlDriverConnectionString = "%s:%s@tcp(%s:%s)/%s"
)

var (
	databaseManager *DatabaseManager
)

type DatabaseManager struct {
	dbConnection *sql.DB
	logManager   log.LogManagerInterface
}

func NewDatabaseManager(logManager log.LogManagerInterface, databaseConfig DatabaseConfig) *DatabaseManager {
	if databaseManager == nil {
		dbConn, err := sql.Open(mysqlDriverName, formatCredentials(databaseConfig))
		if err != nil {
			//todo log
			panic(err)
		}

		databaseManager = &DatabaseManager{dbConnection: dbConn, logManager: logManager}
	}

	return databaseManager
}

func formatCredentials(databaseConfig DatabaseConfig) string {
	return fmt.Sprintf(mysqlDriverConnectionString, databaseConfig.username, databaseConfig.password, databaseConfig.host, databaseConfig.port, databaseConfig.databaseName)
}

func (dm *DatabaseManager) Save(entry types.Entry) error {
	query := "insert into db.passwords (domain, username, password) VALUES (?, ?, ?)"

	stmt, err := dm.dbConnection.Prepare(query)
	if err != nil {
		//todo log
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(entry.Domain, entry.Username, entry.Password)
	if err != nil {
		//todo log
		return err
	}

	//todo log
	return nil
}

func (dm *DatabaseManager) Get(domain string) (*types.Entry, error) {
	query := "SELECT * FROM passwords WHERE domain = ?"

	stmt, err := dm.dbConnection.Prepare(query)
	if err != nil {
		//todo log
		return nil, err
	}

	defer stmt.Close()

	var entry types.Entry

	row := stmt.QueryRow(domain)

	err = row.Scan(&entry.Domain, &entry.Username, &entry.Password)
	if err != nil {
		//todo log
		return nil, err
	}

	//todo log
	return &entry, nil
}

func (dm *DatabaseManager) Contains(domain string) (bool, error) {
	query := "SELECT COUNT(*) FROM passwords WHERE domain = ?"

	stmt, err := dm.dbConnection.Prepare(query)
	if err != nil {
		//todo log
		return false, err
	}

	defer stmt.Close()

	var count int

	err = stmt.QueryRow(domain).Scan(&count)
	if err != nil {
		//todo log
		return false, err
	}

	//todo log
	return count != 0, nil
}

func (dm *DatabaseManager) Update(entry types.Entry) error {
	query := "update db.passwords set username = ?, password = ? where domain = ?"

	stmt, err := dm.dbConnection.Prepare(query)
	if err != nil {
		//todo log
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(entry.Username, entry.Password, entry.Domain)
	if err != nil {
		//todo log
		return err
	}

	//todo log
	return nil
}
