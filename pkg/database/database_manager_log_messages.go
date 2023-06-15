package database

const (
	dbConnectionOpenErrorMessage   = "failed to open db connection: %s"
	queryPreparationFailMessage    = "query preparation failed: %s"
	dbQueryExecFailMessage         = "failed to exec query: %s"
	rowQueryFailMessage            = "failed to query row: %s"
	successfulEntrySaveMessage     = "successfully saved the entry into the db"
	successfulEntryGetMessage      = "successfully get the entry from the db"
	successfulEntryContainsMessage = "successfully check that the db contains the entry"
	successfulEntryUpdateMessage   = "successfully updated the entry in the db"
)
