package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func tableExists(db *sqlx.DB) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='tasks'"
	err := db.Get(&count, query)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func createTable(db *sqlx.DB) error {
	_, err := db.Exec(`CREATE TABLE "tasks" ("id" INTEGER, "title" TEXT NOT NULL, "description" TEXT, "status" TEXT, PRIMARY KEY("id" AUTOINCREMENT))`)
	return err
}

func createTask(db *sqlx.DB, task Task) error {
	_, err := db.Exec("INSERT INTO tasks(title, description, status) VALUES (?,?,?)", task.TaskTitle, task.TaskDescription, task.Status)
	return err
}

func getTasks(db *sqlx.DB) ([]Task, error) {
	var tasks []Task
	_, err := db.Queryx("SELECT * FROM tasks")
	return tasks, err
}

func getTasksByStatus(db *sqlx.DB, status string) ([]Task, error) {
	var tasks []Task
	rows, err := db.Queryx("SELECT * FROM tasks WHERE status=?", status)
	if err != nil {
		return tasks, fmt.Errorf("unable to execute query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.ID,
			&task.TaskTitle,
			&task.TaskDescription,
			&task.Status,
		); err != nil {
			return tasks, fmt.Errorf("unable to scan row: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func getTaskById(db *sqlx.DB, id uint) (Task, error) {
	var task Task
	err := db.QueryRow("SELECT * FROM tasks where id=?", id).Scan(
		&task.ID,
		&task.TaskTitle,
		&task.TaskDescription,
		&task.Status,
	)
	return task, err
}

func updateTask(db *sqlx.DB, task Task) error {
	_, err := db.Exec("UPDATE tasks SET title=?, description=?, status=? WHERE id=?", task.TaskTitle, task.TaskDescription, task.Status, task.ID)
	return err
}

func delTask(db *sqlx.DB, taskID uint) error {
	_, err := db.Exec("DELETE FROM tasks WHERE id=?", taskID)
	return err
}
