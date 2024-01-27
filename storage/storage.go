package storage

import "database/sql"

type PostgresStorage struct {
	db *sql.DB
}

type Task struct {
	Id          int32
	Title       string
	Description string
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (ps *PostgresStorage) AddTask(title, description string) (bool, string) {
	_, err := ps.db.Exec(`INSERT INTO tasks(title, description) VALUES ($1, $2)`, title, description)
	if err != nil {
		return false, "cannot insert"
	}
	return true, "INSERTED SUCCESFULLY"
}

func (ps *PostgresStorage) UpdateTask(task_id int32, title, description string) (bool, string) {
	_, err := ps.db.Exec("UPDATE tasks SET title = $1, description = $2 WHERE id = $3", title, description, task_id)
	if err != nil {
		return false, "cannot update"
	}
	return true, "UPDATED SUCCESFULLY"
}

func (ps *PostgresStorage) DeleteTask(task_id int32) (bool, string) {
	_, err := ps.db.Exec("DELETE from tasks WHERE id = $1", task_id)
	if err != nil {
		return false, "cannot delete"
	}

	return true, "DELETED SUCCESFULLY"
}

func (ps *PostgresStorage) GetOneTask(taskID int32) (int32, string, string, error) {
	var returning_task_id int32
	var title, description string

	if err := ps.db.QueryRow("SELECT id, title, description FROM tasks WHERE id = $1", taskID).Scan(&returning_task_id, &title, &description); err != nil {
		return 0, "", "", err
	}

	return returning_task_id, title, description, nil
}

func (ps *PostgresStorage) GetAllTasks() ([]Task, error) {
	rows, err := ps.db.Query("SELECT id, title, description from tasks")
	if err != nil {
		return nil, err
	}
	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Id, &task.Title, &task.Description); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
