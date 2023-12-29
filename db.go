package main

import (
	"database/sql"
	"reflect"
	"time"
)

type status int

const (
	todo status = iota
	inProgress
	done
)

func (s status) String() string {
	return [...]string{"todo", "inProgress", "done"}[s]
}

type taskDb struct {
	db   *sql.DB
	path string
}

type Task struct {
	Id      int
	Name    string
	Project string
	Status  string
	Created time.Time
}

func (t Task) FilterValue() string {
	return t.Name
}

func (t Task) Title() string {
	return t.Name
}

func (t Task) Description() string {
	return t.Project
}

// check if table exists.
// Maybe this isn't needed because of createTable. Come back to this later.
func (t *taskDb) isTableExists(name string) bool {
	rows, err := t.db.Query("SELECT * FROM " + name)
	if err != nil {
		return false
	}
	defer rows.Close()
	return true
}

// create table if it doesn't exist?
func (t *taskDb) createTable() error {
	_, err := t.db.Exec(`CREATE TABLE IF NOT EXISTS "tasks" 
	( "id" INTEGER, 
		"name" TEXT NOT NULL, 
		"project" TEXT, 
		"status" TEXT, 
		"created" DATETIME,
		 PRIMARY KEY("id" AUTOINCREMENT))`)
	return err
}

func (t *taskDb) insertTask(name, project string) error {
	_, err := t.db.Exec(
		"INSERT INTO tasks(name , project , status,created) VALUES ( ? , ? ,? ,?)",
		name,
		project,
		todo.String(),
		time.Now())
	return err
}

func (t *taskDb) deleteTask(id int) error {
	_, err := t.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}

func (t *taskDb) updateTask(task Task) error {
	origTask, err := t.getTask(task.Id)
	if err != nil {
		return err
	}
	origTask.mergeFields(task)
	_, err = t.db.Exec(
		"UPDATE tasks SET name = ?, project = ?, status = ? WHERE id = ?",
		origTask.Name,
		origTask.Project,
		origTask.Status,
		origTask.Id,
	)
	return err
}

// check for non zero fields and set it in the caller.
func (orig *Task) mergeFields(task Task) {
	newValues := reflect.ValueOf(&task).Elem()
	oldValues := reflect.ValueOf(orig).Elem()

	for i := 0; i < newValues.NumField(); i += 1 {
		newField := newValues.Field(i).Interface()
		if oldValues.CanSet() {
			if v, ok := newField.(int64); ok && newField != 0 {
				oldValues.Field(i).SetInt(v)
			}
			if v, ok := newField.(string); ok && newField != "" {
				oldValues.Field(i).SetString(v)
			}
		}
	}
}

func (t *taskDb) getTasks() ([]Task, error) {
	var tasks []Task
	rows, err := t.db.Query("SELECT * FROM tasks")
	if err != nil {
		return tasks, err
	}

	for rows.Next() {
		var task Task
		err = rows.Scan(
			&task.Id,
			&task.Name,
			&task.Project,
			&task.Status,
			&task.Created,
		)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (t *taskDb) getTask(id int) (Task, error) {
	var task Task
	err := t.db.QueryRow("SELECT * from tasks WHERE id = ?", id).
		Scan(
			&task.Id,
			&task.Name,
			&task.Project,
			&task.Status,
			&task.Created,
		)
	return task, err
}
