package main

import (
	"fmt"
	"gian/cron"
	"gian/db"
	"gian/migrations"
	"gian/router"
	"log"
	"net/http"
	"os"
)

func main() {

	err := db.SetupDB()

	if err != nil {
		log.Fatal("failed to connect database: ", err)
		return
	}

	err = migrations.Migrate()
	if err != nil {
		log.Fatal("failed to migrate: ", err)
		return
	}
	fmt.Println("Hello gaurav")

	cron.RunCron()
	r := router.NewRouter()

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(fmt.Sprintf(":%s", port), r)

}
