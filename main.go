package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tomkeur/mysql-to-strict/checks/date"
	"github.com/tomkeur/mysql-to-strict/checks/datetime"
	"github.com/tomkeur/mysql-to-strict/checks/enum"
	"github.com/tomkeur/mysql-to-strict/database"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/alecthomas/kingpin.v2"
)

const version = "1.4.0"

var (
	dbAddr      = kingpin.Flag("host", "Connect to host. ip:port or hostname:port.").Envar("MYSQL_HOST").Short('h').Required().String()
	dbUsername  = kingpin.Flag("user", "User for login.").Envar("MYSQL_USERNAME").Short('u').Required().String()
	dbPassword  = kingpin.Flag("password", "Password to use when connecting to server. If password is not given it's asked from the tty.").Envar("MYSQL_PASSWORD").Short('p').String()
	dbName      = kingpin.Flag("name", "Database name.").Short('n').Envar("MYSQL_DATABASE").Required().String()
	outputFile  = kingpin.Flag("filename", "Output file").Short('f').Envar("FILENAME").Default("output.sql").String()
	forceUpdate = kingpin.Flag("force", "Force update queries").Bool()
	tables      = Tables{}
	totalTables = 0
	queries     bytes.Buffer
)

type Columns map[string]database.Column
type Tables map[string]Columns

func initializeDatabaseConnection() {
	log.Println("==> Initializing database connection pool...")
	var err error
	// Build data source name.
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", *dbUsername, *dbPassword, *dbAddr, *dbName)
	database.Connection, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := database.Connection.Ping(); err != nil {
		log.Fatal(err)
	}
}

func closeDatabaseConnection() {
	database.Connection.Close()
}

func retrieveAllTables() map[string]string {
	// Get all the MySQL tables for the current database.
	rows, err := database.Connection.Query("SHOW TABLE STATUS;")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	tableNames := make(map[string]string)
	for rows.Next() {
		var t = database.Table{}
		err = rows.Scan(&t.Name, &t.Engine, &t.Version, &t.RowFormat, &t.Rows, &t.AvgRowLength, &t.DataLength, &t.MaxDataLength, &t.IndexLength, &t.DataFree, &t.AutoIncrement, &t.CreateTime, &t.UpdateTime, &t.CheckTime, &t.Collation, &t.CheckSum, &t.CheckTime, &t.CreateOptions)
		if err != nil {
			log.Fatal(err)
		}
		tableNames[t.Name] = t.Engine
	}
	totalTables = len(tableNames)
	log.Println("==> Total tables found: ", totalTables)
	return tableNames
}

func checkTableEngines(tableNames map[string]string) {
	log.Println("==> Starting checking tables engines")
	for tableName, tableEngine := range tableNames {
		if tableEngine != "InnoDB" {
			queries.WriteString(fmt.Sprintf("ALTER TABLE `%s` ENGINE InnoDB;\n", tableName))
		}
	}
}

func describeAllTables(tableNames map[string]string) {
	log.Println("==> Starting checking table structures...")
	for tableName := range tableNames {
		query := fmt.Sprintf("DESCRIBE `%s`;", tableName)
		rows, err := database.Connection.Query(query)
		if err != nil {
			log.Fatal(err)
		}

		columns := Columns{}
		for rows.Next() {
			c := new(database.Column)
			if err := rows.Scan(&c.Field, &c.Type, &c.Null, &c.Key, &c.Default, &c.Extra); err != nil {
				log.Fatal(err)
			}

			// Add all the Columns for this specific table.
			columns[c.Field] = *c
		}

		// Add the current table with columns to the tables.
		tables[tableName] = columns
	}
}

func checkTablesAndFields() {
	log.Println("==> Starting datatypes checks and building queries if needed")
	for tableName, columns := range tables {
		for _, column := range columns {
			// Check if column type is one of the types that is not strict.
			columnType := column.Type.String
			switch columnType {
			case "datetime":
				datetime.Datetime(column, tableName, *forceUpdate)
			case "date":
				date.Date(column, tableName, *forceUpdate)
			case "enum":
				enum.Enum(column, tableName)
			}
		}
	}

	queries.WriteString(datetime.GetQueries())
	queries.WriteString(date.GetQueries())
	queries.WriteString(enum.GetQueries())
}

func main() {
	kingpin.Version(version)
	kingpin.Parse()
	log.Println("==> Starting MySQL to Strict converter ", version)

	// Is password is empty, use readpassword.
	if *dbPassword == "" {
		fmt.Print("Enter MySQL Password:")
		bytePassword, err := terminal.ReadPassword(0)
		if err == nil {
			*dbPassword = string(bytePassword)
		}
	}

	initializeDatabaseConnection()
	tableNames := retrieveAllTables()
	checkTableEngines(tableNames)
	describeAllTables(tableNames)
	checkTablesAndFields()

	if queries.Len() > 0 {
		t := time.Now()
		currentTime := t.Format(time.RFC3339)
		if *outputFile == "output.sql" {
			*outputFile = fmt.Sprintf("%s-%s-%s", *dbName, currentTime, *outputFile)
		}
		file, err := os.Create(fmt.Sprintf("./%s", *outputFile))
		if err != nil {
			panic(err)
		}
		defer file.Close()
		numberOfBytes, err := file.Write(queries.Bytes())
		if err != nil {
			panic(err)
		}
		log.Printf("==> Wrote %d bytes to the file: %s", numberOfBytes, *outputFile)
	} else {
		log.Println("==> No work needed!")
	}

	// Close the MySQL connection.
	defer closeDatabaseConnection()
	log.Println("==> This is the end my friend, have a nice day!")
}
