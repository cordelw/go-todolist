package controller

import (
	"gotodo/models"
	"net/http"
)

// ROUTES //
// Product page
func ServeIndexPage(w http.ResponseWriter, r *http.Request) {
	RenderHTMLTemplate("index", "./pages/index.html", w, nil)
}

// Web app
func (c *Controller) ServeAppPage(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userid").(string)
	userName := r.Context().Value("username").(string)
	data := struct {
		Username string
		Tasks    []models.Task
		Finished []models.Task
	}{userName, models.TaskQueryIncompleteByUserID(c.DB, userId), models.TaskQueryCompleteByUserID(c.DB, userId)}

	RenderHTMLTemplate("app", "./pages/app.html", w, data)
}

func (c *Controller) ServeRegisterPage(w http.ResponseWriter, r *http.Request) {
	RenderHTMLTemplate("register", "./pages/register.html", w, nil)
}

func (c *Controller) ServeLoginPage(w http.ResponseWriter, r *http.Request) {
	RenderHTMLTemplate("login", "./pages/login.html", w, nil)
}
