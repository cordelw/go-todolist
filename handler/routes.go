package handler

import (
	"database/sql"
	"gotodo/tables"
	"html/template"
	"net/http"
)

// Context
type Context struct {
	DB *sql.DB
}

func NewContext(db *sql.DB) Context {
	return Context{
		DB: db,
	}
}

// Routes
func (ctx Context) ServeIndexPage(w http.ResponseWriter, r *http.Request) {
	data := struct{ Tasks []tables.Task }{tables.TaskQueryByIdDesc(ctx.DB)}
	tmpl, _ := template.New("index").ParseFiles("./pages/index.html")
	tmpl.Execute(w, data)
}
