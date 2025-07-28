package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Todo represents a todo item
type Todo struct {
	ID        int64
	Title     string
	Done      bool
	CreatedAt string
	UpdatedAt string
}

// TodoStorage interface defines methods for Todo storage operations
type TodoStorageInterface interface {
	Save(todo *Todo) error
	Update(todo *Todo) error
	GetByID(id int64) (*Todo, error)
	GetAll() ([]Todo, error)
	GetOldestUndone() (*Todo, error)
	Delete(id int64) error
}

// SQLiteTodoStorage implements TodoStorage interface
type SQLiteTodoStorage struct {
	db *sql.DB
}

// NewSQLiteTodoStorage creates a new SQLite storage instance
func NewSQLiteTodoStorage(dbPath string) (*SQLiteTodoStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			done BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, err
	}

	return &SQLiteTodoStorage{db: db}, nil
}

// Save stores a new todo item
func (s *SQLiteTodoStorage) Save(todo *Todo) error {
	result, err := s.db.Exec(
		"INSERT INTO todos (title, done) VALUES (?, ?)",
		todo.Title, todo.Done,
	)
	if err != nil {
		return err
	}
	todo.ID, _ = result.LastInsertId()
	return nil
}

// Update modifies an existing todo item
func (s *SQLiteTodoStorage) Update(todo *Todo) error {
	_, err := s.db.Exec(
		"UPDATE todos SET title=?, done=?, updated_at=CURRENT_TIMESTAMP WHERE id=?",
		todo.Title, todo.Done, todo.ID,
	)
	return err
}

// GetByID retrieves a todo item by its ID
func (s *SQLiteTodoStorage) GetByID(id int64) (*Todo, error) {
	todo := &Todo{}
	err := s.db.QueryRow(
		"SELECT id, title, done, created_at, updated_at FROM todos WHERE id=?",
		id,
	).Scan(&todo.ID, &todo.Title, &todo.Done, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

// GetAll retrieves all todo items
func (s *SQLiteTodoStorage) GetAll() ([]Todo, error) {
	rows, err := s.db.Query("SELECT id, title, done, created_at, updated_at FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Done, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

// GetOldestUndone retrieves the oldest todo item that is not done
func (s *SQLiteTodoStorage) GetOldestUndone() (*Todo, error) {
	todo := &Todo{}
	err := s.db.QueryRow(
		"SELECT id, title, done, created_at, updated_at FROM todos WHERE done = FALSE ORDER BY created_at ASC LIMIT 1",
	).Scan(&todo.ID, &todo.Title, &todo.Done, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

// Delete removes a todo item by its ID
func (s *SQLiteTodoStorage) Delete(id int64) error {
	_, err := s.db.Exec("DELETE FROM todos WHERE id=?", id)
	return err
}
