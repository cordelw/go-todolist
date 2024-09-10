package main

import (
	"database/sql"
	"html/template"
	"io"
	"net/http"

	_ "modernc.org/sqlite"

	"github.com/labstack/echo/v4"
)

// Templating
//

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplRenderer() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
}

// Database Stuff
//

type Task struct {
	Id   int
	Text string
}

func dbQueryAll(db *sql.DB) []Task {
	var tasks []Task

	rows, err := db.Query("SELECT id, text FROM tasks")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var task Task
		rows.Scan(&task.Id, &task.Text)

		tasks = append(tasks, task)
	}

	return tasks
}

func dbQueryByOrigin(db *sql.DB, origin string) []Task {
	var tasks []Task

	rows, err := db.Query("SELECT id, task FROM tasks WHERE origin=?", origin)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var task Task
		rows.Scan(&task.Id, &task.Text)

		tasks = append(tasks, task)
	}

	return tasks
}

// main
//

const listenPort string = ":8080"

type Page struct {
	Tasks []Task
}

func newPage() Page {
	return Page{
		Tasks: []Task{},
	}
}

func main() {
	const dsn string = "database.db"
	const dbSchema = `
		CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		origin TEXT NOT NULL,
		task TEXT NOT NULL
		);`

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if _, err := db.Exec(dbSchema); err != nil {
		panic(err)
	}

	ctx := echo.New()
	ctx.Renderer = newTemplRenderer()

	data := newPage()

	ctx.GET("/", func(c echo.Context) error {
		origin := c.RealIP()

		data.Tasks = dbQueryByOrigin(db, origin)
		return c.Render(http.StatusOK, "index", data)
	})

	ctx.POST("/create-task", func(c echo.Context) error {
		origin := c.RealIP()
		taskInput := c.FormValue("task-input")

		if taskInput == "" {
			return c.Render(http.StatusBadRequest, "body", nil)
		}

		_, err := db.Exec(`INSERT INTO tasks VALUES(NULL, ?, ?)`, origin, taskInput)
		if err != nil {
			return c.Render(http.StatusInternalServerError, "body", nil)
		}

		data.Tasks = dbQueryByOrigin(db, origin)
		return c.Render(http.StatusOK, "index", data)
	})

	ctx.POST("/delete-task", func(c echo.Context) error {
		origin := c.RealIP()
		taskID := c.FormValue("task-id")

		var realorigin string
		db.QueryRow("SELECT origin FROM tasks WHERE ID=?", taskID).Scan(&realorigin)

		if origin != realorigin {
			return c.Render(http.StatusBadRequest, "body", nil)
		}

		db.Exec("DELETE FROM tasks WHERE id=?", taskID)

		data.Tasks = dbQueryByOrigin(db, origin)
		return c.Render(http.StatusOK, "task-list", data)
	})

	// Start server
	ctx.Logger.Fatal(ctx.Start(listenPort))
}
