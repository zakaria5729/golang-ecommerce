package user

import "database/sql"

func FetchUser(db *sql.DB, id int) (string, error) {
	// Example query
	return "DB access for user " + string(rune(id)), nil
}
