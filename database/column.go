package database

import "database/sql"

type Column struct {
	Field   string
	Type    sql.NullString
	Null    sql.NullString
	Key     sql.NullString
	Default sql.NullString
	Extra   sql.NullString
}
