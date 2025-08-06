package tables

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

const TaskSchema = `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT,
		completed INTEGER);`

type Task struct {
	Id   int
	Text string
}

func TaskQueryByIdDesc(db *sql.DB) []Task {
	var tasks []Task

	rows, err := db.Query("SELECT id, text FROM tasks WHERE completed=false ORDER BY id DESC")
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

func TaskCreate(db *sql.DB, text string) error {
	_, err := db.Exec("INSERT INTO tasks (text, completed) VALUES (?, false)", text)
	return err
}

func TaskComplete(db *sql.DB, taskID string) error {
	_, err := db.Exec("UPDATE tasks SET completed='1' WHERE id=?", taskID)
	return err
}
