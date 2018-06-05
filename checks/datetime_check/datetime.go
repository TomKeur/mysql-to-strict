package datetime_check

import (
	"strings"
	"fmt"
	"bytes"
	"github.com/tomkeur/mysql-to-strict/database"
)

var queries bytes.Buffer

func Datetime(column database.Column, tableName string) {
	// Checks.
	defaultTimeStamp := strings.Contains(column.Default.String, "0000-00-00 00:00:00")
	defaultNull := column.Null.String == "NO"
	if defaultTimeStamp || defaultNull {
		buildQuery(column, tableName)
	}
}

func buildQuery(column database.Column, tableName string) {
	alterQuery := "" +
		"ALTER TABLE `%s`\n" +
		"	CHANGE COLUMN `%s` `%s` DATETIME NULL DEFAULT NULL;\n"
	updateDataQuery := "" +
		"UPDATE `%s` SET `%s` = NULL WHERE `%s` = '0000-00-00 00:00:00';\n"
	field := column.Field
	queries.WriteString(fmt.Sprintf(alterQuery, tableName, field, field))
	queries.WriteString(fmt.Sprintf(updateDataQuery, tableName, field, field))
}

func GetQueries() string {
	return queries.String()
}
