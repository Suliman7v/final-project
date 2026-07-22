package tests

import (
	"log"
	"os"
	"testing"

	"final-project/pkg/db"
)

func TestMain(m *testing.M) {
	err := db.Init(TestDBConnStr)
	if err != nil {
		log.Fatal("Не удалось подключиться к тестовой БД:", err)
	}

	_, err = db.DB.Exec(`
		DROP TABLE IF EXISTS scheduler CASCADE;
		CREATE TABLE scheduler (
			id SERIAL PRIMARY KEY,
			date CHAR(8) NOT NULL DEFAULT '',
			title VARCHAR(255) NOT NULL DEFAULT '',
			comment TEXT NOT NULL DEFAULT '',
			repeat VARCHAR(128) NOT NULL DEFAULT ''
		);
		CREATE INDEX idx_date ON scheduler(date);
	`)
	if err != nil {
		log.Fatal("Ошибка при создании схемы:", err)
	}

	code := m.Run()

	db.Close()

	os.Remove("../scheduler.db")

	os.Exit(code)
}
