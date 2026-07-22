package db

import (
	"database/sql"
	"fmt"
	"strconv"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (string, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) 
	          VALUES ($1, $2, $3, $4) 
	          RETURNING id`

	var id int
	err := DB.QueryRow(query, task.Date, task.Title, task.Comment, task.Repeat).Scan(&id)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(int64(id), 10), nil
}

func Tasks(limit int, search string) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT $1`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task

	for rows.Next() {
		task := &Task{}
		var id int
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		task.ID = strconv.Itoa(id)
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

func GetTask(idStr string) (*Task, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = $1`
	row := DB.QueryRow(query, id)

	task := &Task{}
	var dbID int
	err = row.Scan(&dbID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	task.ID = strconv.Itoa(dbID)
	return task, nil
}

func UpdateTask(task *Task) error {
	id, err := strconv.Atoi(task.ID)
	if err != nil {
		return err
	}

	query := `UPDATE scheduler SET date = $1, title = $2, comment = $3, repeat = $4 WHERE id = $5`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

func DeleteTask(idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	query := `DELETE FROM scheduler WHERE id = $1`
	res, err := DB.Exec(query, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

func UpdateDate(idStr string, date string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	query := `UPDATE scheduler SET date = $1 WHERE id = $2`
	res, err := DB.Exec(query, date, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}
