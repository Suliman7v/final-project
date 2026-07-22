package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"final-project/pkg/db"
)

type Task struct {
	ID      int64  `db:"id"`
	Date    string `db:"date"`
	Title   string `db:"title"`
	Comment string `db:"comment"`
	Repeat  string `db:"repeat"`
}

func countTasks() (int, error) {
	var count int
	err := db.DB.QueryRow("SELECT COUNT(id) FROM scheduler").Scan(&count)
	return count, err
}

func TestDB(t *testing.T) {
	InitTestDB(t)
	defer CloseTestDB(t)

	cleanupDB(t)

	before, err := countTasks()
	assert.NoError(t, err)

	today := time.Now().Format("20060102")

	var id int
	err = db.DB.QueryRow(`
		INSERT INTO scheduler (date, title, comment, repeat) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id`,
		today, "Todo", "Комментарий", "").Scan(&id)
	assert.NoError(t, err)

	var task struct {
		ID      int
		Date    string
		Title   string
		Comment string
		Repeat  string
	}
	err = db.DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = $1", id).
		Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	assert.NoError(t, err)
	assert.Equal(t, id, task.ID)
	assert.Equal(t, "Todo", task.Title)
	assert.Equal(t, "Комментарий", task.Comment)

	_, err = db.DB.Exec("DELETE FROM scheduler WHERE id = $1", id)
	assert.NoError(t, err)

	after, err := countTasks()
	assert.NoError(t, err)

	assert.Equal(t, before, after)
}
