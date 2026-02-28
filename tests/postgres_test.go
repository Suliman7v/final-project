package tests

import (
	"log"
	"os"
	"testing"

	"final-project/pkg/db"
)

func TestMain(m *testing.M) {
	// Подключаемся к тестовой PostgreSQL
	connStr := "postgres://postgres:password@localhost:5432/todo_test?sslmode=disable"

	err := db.Init(connStr)
	if err != nil {
		log.Fatal("Не удалось подключиться к тестовой БД:", err)
	}

	// Очищаем и создаем таблицу
	db.DB.Exec("DROP TABLE IF EXISTS scheduler")

	schema := `
	CREATE TABLE scheduler (
		id SERIAL PRIMARY KEY,
		date CHAR(8) NOT NULL DEFAULT '',
		title VARCHAR(255) NOT NULL DEFAULT '',
		comment TEXT NOT NULL DEFAULT '',
		repeat VARCHAR(128) NOT NULL DEFAULT ''
	);
	CREATE INDEX idx_date ON scheduler(date);
	`

	_, err = db.DB.Exec(schema)
	if err != nil {
		log.Fatal("Ошибка при создании схемы:", err)
	}

	os.Exit(m.Run())
}
