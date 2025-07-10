package database

import (
	"database/sql"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

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

func Tasks(limit int, search string) ([]*Task, error) {
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
	} else {
		query += "ORDER BY date LIMIT :limit"
		rows, err = DB.Query(query, sql.Named("limit", limit))
	}
	if err != nil {
		return make([]*Task, 0), err
	}
	defer rows.Close()
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return make([]*Task, 0), err
		}
		tasks = append(tasks, &t)

	}
	if err := rows.Err(); err != nil {
		return make([]*Task, 0), err
	}
	if tasks == nil {
		tasks = make([]*Task, 0)
	}
	return tasks, nil
}
