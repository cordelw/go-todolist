package controller

import (
	"database/sql"
	"net/http"
	"text/template"
)

func RenderHTMLTemplate(templateName string, templatePath string, w http.ResponseWriter, data any) {
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error"))
		return
	}

	t.ExecuteTemplate(w, templateName, data)
}

// Context
type Controller struct {
	DB *sql.DB
}

func NewController(db *sql.DB) Controller {
	return Controller{
		DB: db,
	}
}
