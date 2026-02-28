package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func notFoundTask(t *testing.T, id string) {
	body, err := requestJSON("api/task?id="+id, nil, http.MethodGet)
	assert.NoError(t, err)
	var m map[string]any
	err = json.Unmarshal(body, &m)
	assert.NoError(t, err)
	_, ok := m["error"]
	assert.True(t, ok)
}

func TestDone(t *testing.T) {
	cleanupDB(t)

	// Создаем задачу с повторением
	task := map[string]interface{}{
		"date":    "20260228",
		"title":   "Фитнес",
		"comment": "с тренером",
		"repeat":  "d 3",
	}

	id := CreateTask(t, task)

	// Отмечаем выполненной
	resp, err := http.Post("http://localhost:7540/api/task/done?id="+id, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", resp.StatusCode)
	}

	// Проверяем, что дата изменилась (должна стать +3 дня)
	updated := GetTask(t, id)

	expectedDate := "20260303" // 28.02 + 3 дня = 03.03
	if updated["date"] != expectedDate {
		t.Errorf("Дата не обновилась: ожидалось %s, получено %s", expectedDate, updated["date"])
	}
}

func TestDelTask(t *testing.T) {
	db := openDB(t)
	defer db.Close()

	id := addTask(t, task{
		title:  "Временная задача",
		repeat: "d 3",
	})
	ret, err := postJSON("api/task?id="+id, nil, http.MethodDelete)
	assert.NoError(t, err)
	assert.Empty(t, ret)

	notFoundTask(t, id)

	ret, err = postJSON("api/task", nil, http.MethodDelete)
	assert.NoError(t, err)
	assert.NotEmpty(t, ret)
	ret, err = postJSON("api/task?id=wjhgese", nil, http.MethodDelete)
	assert.NoError(t, err)
	assert.NotEmpty(t, ret)
}
