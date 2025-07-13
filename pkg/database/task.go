package database

import (
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Функция для реализации INSERT запроса в базу данных, возвращает id добавленной задачи
func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment,repeat) VALUES (:date, :title, :com, :repeat)`
	res, err := DB.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("com", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

// Функция для реализации SELECT запроса в базу данных по указанному id
func GetTask(id string) (*Task, error) {
	var task Task
	err := DB.QueryRow("SELECT * from scheduler WHERE id = :id", sql.Named("id", id)).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// Функция для реализации UPDATE запроса в базу данных по указанному id
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat= ? WHERE id = ?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

// Функция для реализации DELETE запроса в базу данных по указанному id
func DeleteTask(id string) error {
	res, err := DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for deleting task`)
	}
	return nil
}

// Функция для реализации SELECT запроса в базу данных, для поиска задач по параметру search
func GetTasks(limit int, search string) ([]*Task, error) {
	var tasks []*Task
	query := "SELECT * from scheduler "
	var whereQuery string
	var rows *sql.Rows
	var err error
	if search != "" {
		date, err := time.Parse("02.01.2006", search)
		if err != nil {
			search = "%" + search + "%"
			whereQuery = "WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit"
		} else {
			search = date.Format("20060102")
			whereQuery = "WHERE date LIKE :search LIMIT :limit"
		}
		query = query + whereQuery
		rows, err = DB.Query(query, sql.Named("limit", limit), sql.Named("search", search))
		if err != nil {
			return nil, err
		}
	} else {
		query += "ORDER BY date LIMIT :limit"
		rows, err = DB.Query(query, sql.Named("limit", limit))
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if tasks == nil {
		tasks = make([]*Task, 0)
	}
	return tasks, nil
}
