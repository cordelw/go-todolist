package handler

import (
	"gotodo/tables"
	"net/http"
)

func (c Context) ApiCreateTask(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Could not parse form.", http.StatusBadRequest)
		return
	}

	taskstr := r.PostFormValue("task-input")
	err := tables.TaskCreate(c.DB, taskstr)
	if err != nil {
		http.Error(w, "Could not create task", http.StatusBadRequest)
	}

	c.ServeIndexPage(w, r)
}

func (c Context) ApiCompleteTask(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Could not parse form.", http.StatusBadRequest)
		return
	}

	taskID := r.PostFormValue("task-id")
	if err := tables.TaskComplete(c.DB, taskID); err != nil {
		http.Error(w, "Could not complete task.", http.StatusBadRequest)
	}

	c.ServeIndexPage(w, r)
}
