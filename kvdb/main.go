package main

import (
	"log"
	"os"

	"github.com/kxrxh/kvdb/api"
	"github.com/kxrxh/kvdb/storage"
)

func main() {
	// Init database
	dbFile, exists := os.LookupEnv("DB_FILE")
	if !exists {
		dbFile = "db.json"
	}

	err := storage.NewKeyValueStore(dbFile)
	if err != nil {
		panic(err)
	}

	// Creating API
	app := api.InitApp()
	api.RegisterRoutes(app)

	dbPort, exist := os.LookupEnv("DB_PORT")
	if !exist {
		dbPort = "5173"
	}

	// Start API
	log.Printf("Starting API on port: %s", dbPort)
	app.Listen(":" + dbPort)
}
