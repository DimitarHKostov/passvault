package database

var (
	databaseManager *DatabaseManager
)

type DatabaseManager struct{}

func Get() *DatabaseManager {
	if databaseManager == nil {
		databaseManager = &DatabaseManager{}
	}

	return databaseManager
}

func (dm *DatabaseManager) Save() {

}

