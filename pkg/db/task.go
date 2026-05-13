package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

const DateFormat = "20060102"
const DefaultLimit = 50

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat)
	          VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

func GetTasks(limit int, search string) ([]Task, error) {
	if limit <= 0 {
		limit = DefaultLimit
	}
	query, args := buildQuery(limit, search)
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		var id int64
		if err := rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		t.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	if tasks == nil {
		tasks = []Task{}
	}
	return tasks, nil
}

func buildQuery(limit int, search string) (string, []interface{}) {
	if search == "" {
		return "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?",
			[]interface{}{limit}
	}

	if t, err := time.Parse("02.01.2006", search); err == nil {
		dateStr := t.Format(DateFormat)
		return "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?",
			[]interface{}{dateStr, limit}
	}

	like := "%" + search + "%"
	return "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?",
		[]interface{}{like, like, limit}
}

func GetTask(id string) (*Task, error) {
	var t Task
	var idInt int64
	err := DB.QueryRow(
		"SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?",
		id,
	).Scan(&idInt, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("задача не найдена")
		}
		return nil, err
	}
	t.ID = strconv.FormatInt(idInt, 10)
	return &t, nil
}

func UpdateTask(task *Task) error {
	res, err := DB.Exec(
		"UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		task.Date, task.Title, task.Comment, task.Repeat, task.ID,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func DeleteTask(id string) error {
	res, err := DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func UpdateTaskDate(id string, date string) error {
	res, err := DB.Exec("UPDATE scheduler SET date = ? WHERE id = ?", date, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}