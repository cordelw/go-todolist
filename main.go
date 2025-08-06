package main

import (
	"database/sql"
	"gotodo/handler"
	"gotodo/middleware"
	"gotodo/tables"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

const lport string = ":8080"
const dsn string = "data.db"

func main() {
	// Connect to database
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Execute database schemas
	if _, err := db.Exec(tables.TaskSchema); err != nil {
		log.Fatal(err)
	}

	// Router + routes
	router := http.NewServeMux()
	ctx := handler.NewContext(db)

	router.HandleFunc("GET /", ctx.ServeIndexPage)
	router.HandleFunc("POST /create-task", ctx.ApiCreateTask)
	router.HandleFunc("POST /complete-task", ctx.ApiCompleteTask)

	// HTTP server
	server := http.Server{
		Addr:    lport,
		Handler: middleware.Logging(router),
	}

	// Listen and serve
	if err = server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
