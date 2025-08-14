package models

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

const TasksSchema = `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		userid INTEGER SECONDAY KEY NOT NULL,
		text TEXT NOT NULL,
		completed INTEGER);`

type Task struct {
	Id       int
	UserId   int
	Text     string
	Complete bool
}

func TaskQueryByUserID(db *sql.DB, userId string) []Task {
	var tasks []Task

	rows, err := db.Query("SELECT id, userid, text FROM tasks WHERE completed=false AND userid=? ORDER BY id DESC", userId)
	if err != nil {
		return tasks
	}

	for rows.Next() {
		var task Task
		rows.Scan(&task.Id, &task.UserId, &task.Text)

		tasks = append(tasks, task)
	}

	return tasks
}

func TaskQueryByID(db *sql.DB, id string) Task {
	var task Task

	row := db.QueryRow("SELECT id, userid, text FROM tasks WHERE id=?", id)
	if err := row.Scan(&task.Id, &task.UserId, &task.Text); err != nil {
		return Task{}
	}

	return task
}

func TaskCreate(db *sql.DB, userId string, text string) error {
	_, err := db.Exec("INSERT INTO tasks (userid, text, completed) VALUES (?, ?, false)", userId, text)
	return err
}

func TaskComplete(db *sql.DB, taskId string) error {
	_, err := db.Exec("UPDATE tasks SET completed='1' WHERE id=?", taskId)
	return err
}
