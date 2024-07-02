package utils

const (
	ErrNoUsersFound = `no users found`

	ErrQuery   = "failed to execute query: %v with params: %+v, error: %v"
	ErrScanRow = "failed to scan row for query: %v with params: %+v, error: %v"

	ErrRowIteration = "row iteration error for query: %v with params: %+v, error: %v"

	ErrBeginTx = "cant begin tx: %v"
)

const (
	ErrCantOpenDB = "cant open db"
	ErrCantPingDB = "cant ping db"
)
