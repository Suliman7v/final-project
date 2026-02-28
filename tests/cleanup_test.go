package tests

import (
	"testing"

	"final-project/pkg/db"
)

func cleanupDB(t *testing.T) {
	_, err := db.DB.Exec("TRUNCATE scheduler RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatal("Не удалось очистить БД:", err)
	}
}
