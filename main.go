package main

import (
	"database/sql"
	"gotodo/controller"
	"gotodo/middleware"
	"gotodo/models"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/adrianosela/sslmgr"

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
	if _, err := db.Exec(models.TasksSchema); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(models.UsersSchema); err != nil {
		log.Fatal(err)
	}

	// Router + routes
	router := http.NewServeMux()
	staticfs := http.FileServer(http.Dir("./static"))
	router.Handle("GET /static/", http.StripPrefix("/static/", staticfs))
	c := controller.NewController(db)

	// Static pages
	router.HandleFunc("GET /", controller.ServeIndexPage)
	router.HandleFunc("GET /faq", controller.ServeFaqPage)

	router.HandleFunc("GET /register", c.ServeRegisterPage)
	router.HandleFunc("GET /login", c.ServeLoginPage)

	// API
	// auth
	router.HandleFunc("POST /api/register-user", c.ApiUserRegister)
	router.HandleFunc("POST /api/login-user", c.ApiUserLogin)
	router.HandleFunc("POST /api/logout-user", c.ApiUserLogout)

	// tasks
	router.HandleFunc("POST /api/create-task", c.RequiresAuth(c.ApiTaskCreate))
	router.HandleFunc("POST /api/complete-task", c.RequiresAuth(c.ApiTaskComplete))
	router.HandleFunc("POST /api/uncomplete-task", c.RequiresAuth(c.ApiTaskUncomplete))
	router.HandleFunc("POST /api/update-task", c.RequiresAuth(c.ApiTaskUpdate))
	router.HandleFunc("GET /app", c.RequiresAuth(c.ServeAppPage))

	// HTTP server
	/*server := http.Server{
		Addr:    lport,
		Handler: middleware.Logging(router),
	}*/

	// Listen and serve
	/*if err = server.ListenAndServe(); err != nil {
	log.Fatal(err)
	}*/

	server, err := sslmgr.NewServer(sslmgr.ServerConfig{
		Hostnames: []string{"cordelw.com", "www.cordelw.com"},
		HTTPPort:  ":8080",
		HTTPSPort: ":443",
		Handler:   middleware.Logging(router),
		ServeSSLFunc: func() bool {
			return strings.ToLower(os.Getenv("DEV")) != "true"
		},
	})
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	server.ListenAndServe()
}
