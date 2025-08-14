package controller

import (
	"gotodo/models"
	"net/http"
	"strconv"
)

func (c Controller) ApiTaskCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Could not parse form.", http.StatusBadRequest)
		return
	}

	// user-id from request context
	userId := r.Context().Value("userid").(string)

	// Manipulate DB
	taskstr := r.PostFormValue("task-input")
	err := models.TaskCreate(c.DB, userId, taskstr)
	if err != nil {
		http.Error(w, "Could not create task", http.StatusBadRequest)
		return
	}

	c.ServeAppPage(w, r)
}

func (c Controller) ApiTaskComplete(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Could not parse form.", http.StatusBadRequest)
		return
	}
	taskId := r.PostFormValue("task-id")

	// Check if logged in userid matches task's userid
	userId := r.Context().Value("userid").(string)
	if task := models.TaskQueryByID(c.DB, taskId); strconv.Itoa(task.UserId) != userId {
		http.Error(w, "Invalid operation.", http.StatusForbidden)
		return
	}

	// Manipulate db
	if err := models.TaskComplete(c.DB, taskId); err != nil {
		http.Error(w, "Could not complete task.", http.StatusBadRequest)
	}

	c.ServeAppPage(w, r)
}
