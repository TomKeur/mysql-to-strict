package date

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/tomkeur/mysql-to-strict/database"
)

var queries bytes.Buffer

func Date(column database.Column, tableName string, forceUpdate bool) {
	// Checks.
	defaultTimeStamp := strings.Contains(column.Default.String, "0000-00-00")
	defaultNull := column.Null.String == "NO"
	if defaultTimeStamp || defaultNull || forceUpdate {
		buildQuery(column, tableName)
	}
}

func buildQuery(column database.Column, tableName string) {
	alterQuery := "" +
		"ALTER TABLE `%s`\n" +
		"	CHANGE COLUMN `%s` `%s` DATE NULL DEFAULT NULL;\n"
	updateDataQuery := "" +
		"UPDATE `%s` SET `%s` = NULL WHERE `%s` = '0000-00-00';\n"
	field := column.Field
	queries.WriteString(fmt.Sprintf(alterQuery, tableName, field, field))
	queries.WriteString(fmt.Sprintf(updateDataQuery, tableName, field, field))
}

func GetQueries() string {
	return queries.String()
}
