package main

import (
	"app/api"
	database_postgres "app/infrascture/database/postgres"
)

func main() {
	println("Starting application...")

	// Attempt internal connection logic (will fail without real env vars, but tests compilation)
	// oracledb.InitConnection(cfg)
	// Commented out actual connection to avoid panic during simple build/run test without env vars
	// You can uncomment it to test actual connection if env vars are present.

	// For now, just verifying import and signature matches
	db := database_postgres.ConnectDB()
	database_postgres.RunMigrations(db)

	r := api.SetupRouter(db)
	r.Run(":8080")
}
