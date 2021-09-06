package postgres

import (
	"database/sql"
	"os"

	"github.com/Todorov99/sensorcli/pkg/logger"
	_ "github.com/lib/pq"
)

// DatabaseConnection opens postgres connection.
var DatabaseConnection *sql.DB

var postgresLogger = logger.NewLogrus("postres", os.Stdout)

func init() {
	postgresLogger.Info("Initializing postgres DB client")
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
