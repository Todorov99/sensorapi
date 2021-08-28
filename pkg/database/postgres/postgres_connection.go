package postgres

import (
	"database/sql"

	"github.com/Todorov99/sensorcli/pkg/logger"
	_ "github.com/lib/pq"
)

// DatabaseConnection opens postgres connection.
var DatabaseConnection *sql.DB

var postgresLogger = logger.NewLogger("./postres")

func init() {
	db, err := sql.Open("postgres", "postgres://postgres:Password321@localhost/sensorCLI?sslmode=disable")
	if err != nil {
		postgresLogger.Panic(err)
	}

	err = db.Ping()
	if err != nil {
		postgresLogger.Panic(err)
	}

	DatabaseConnection = db
}
