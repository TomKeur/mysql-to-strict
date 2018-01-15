package enum_check

import (
	"bytes"
	"strings"
	"fmt"
	"github.com/tomkeur/mysql-to-strict/database"
)

var queries bytes.Buffer

func Enum(column database.Column, tableName string) {
	// Checks.
	emptyValueInEnums := strings.Contains(column.Type.String, "''")
	defaultNull := column.Null.String == "YES"
	if defaultNull {
		buildNotNullQuery(column, tableName)
	}
	if emptyValueInEnums {
		buildQuery(column, tableName)
	}
}

func buildNotNullQuery(column database.Column, tableName string) {
	runes := []rune(column.Type.String)
	enumString := string(runes[5:(len(column.Type.String) - 1)])

	enumValues := strings.Split(enumString, ",")
	alterQuery := "" +
		"ALTER TABLE `%s`\n" +
		"	CHANGE COLUMN `%s` `%s` ENUM(%s) NOT NULL;\n"
	field := column.Field
	queries.WriteString(fmt.Sprintf(alterQuery, tableName, field, field, strings.Join(enumValues, ",")))
}

func buildQuery(column database.Column, tableName string) {
	runes := []rune(column.Type.String)
	enumString := string(runes[5:(len(column.Type.String) - 1)])

	// Check if there are empy values in the ENUM field.
	enumValues := strings.Split(enumString, ",")
	runAlterQuery := false
	for key, value := range enumValues {
		if value == "''" {
			enumValues[key] = "'UNKNOWN'"
			runAlterQuery = true
		}
	}
	alterQuery := "" +
		"ALTER TABLE `%s`\n" +
		"	CHANGE COLUMN `%s` `%s` ENUM(%s) NOT NULL DEFAULT 'UNKNOWN';\n"

	updateDataQuery := "" +
		"UPDATE `%s` SET `%s` = 'UNKNOWN' WHERE `%s` = '';\n"
	field := column.Field

	// Check if we need to update the empty enum values to the newly created UNKOWN field.
	if runAlterQuery {
		queries.WriteString(fmt.Sprintf(alterQuery, tableName, field, field, strings.Join(enumValues, ",")))
	}
	queries.WriteString(fmt.Sprintf(updateDataQuery, tableName, field, field))
}

func GetQueries() string {
	return queries.String()
}
