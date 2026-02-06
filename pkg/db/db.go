package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
    );
CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`

var DB *sql.DB

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err) // Правильная проверка

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	DB = db

	if install {
		_, err = DB.Exec(schema)
		if err != nil {
			return err
		}
	}

	return db.Ping()
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
