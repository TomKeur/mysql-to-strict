package database

import "database/sql"

type Table struct {
	Name          string
	Engine        string
	Version       sql.NullInt64
	RowFormat     sql.NullString
	Rows          sql.NullInt64
	AvgRowLength  sql.NullInt64
	DataLength    sql.NullInt64
	MaxDataLength sql.NullInt64
	IndexLength   sql.NullInt64
	DataFree      sql.NullInt64
	AutoIncrement sql.NullInt64
	CreateTime    sql.NullString
	UpdateTime    sql.NullString
	CheckTime     sql.NullString
	Collation     sql.NullString
	CheckSum      sql.NullString
	CreateOptions sql.NullString
	Columns       map[string]Column
}
