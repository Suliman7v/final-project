package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id SERIAL PRIMARY KEY,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`

var DB *sql.DB

func Init(connStr string) error {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	DB = db

	if err = db.Ping(); err != nil {
		return err
	}

	_, err = DB.Exec(schema)
	return err
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
