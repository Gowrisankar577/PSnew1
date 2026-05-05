package main

import (
	"log"

	"ps_portal/config"
	"ps_portal/db"
	"ps_portal/routes"
)

func main() {
	// 🔗 Initialize DB
	db.InitDB()
	if db.DB == nil {
		log.Fatal("❌ DB not initialized properly")
	}
	defer func() {
		if err := db.DB.Close(); err != nil {
			log.Println("⚠️  Error closing DB:", err)
		}
	}()

	// ⚙️ Load .env config, Redis, etc.
	config.LoadConfig()

	if err := routes.SetupRouter().Run(":8080"); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
