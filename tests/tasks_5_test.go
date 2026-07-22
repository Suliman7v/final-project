package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func addTask(t *testing.T, task task) string {
	ret, err := postJSON("api/task", map[string]any{
		"date":    task.date,
		"title":   task.title,
		"comment": task.comment,
		"repeat":  task.repeat,
	}, http.MethodPost)
	assert.NoError(t, err)
	assert.NotNil(t, ret["id"])
	id := fmt.Sprint(ret["id"])
	assert.NotEmpty(t, id)
	return id
}

func getTasks(t *testing.T) []map[string]string {
	body, err := requestJSON("api/tasks", nil, http.MethodGet)
	assert.NoError(t, err)

	var m map[string][]map[string]string
	err = json.Unmarshal(body, &m)
	assert.NoError(t, err)
	return m["tasks"]
}

func TestTasks(t *testing.T) {
	cleanupDB(t)

	tasks := []map[string]interface{}{
		{
			"date":    "20260228",
			"title":   "Фитнес",
			"comment": "с тренером",
			"repeat":  "d 1",
		},
		{
			"date":    "20260301",
			"title":   "Сходить в бассейн",
			"comment": "с друзьями",
			"repeat":  "",
		},
		{
			"date":    "20260305",
			"title":   "Уроки",
			"comment": "математика",
			"repeat":  "d 7",
		},
	}

	var ids []string
	for _, task := range tasks {
		id := CreateTask(t, task)
		ids = append(ids, id)
	}

	tasksList := getTasks(t)

	if len(tasksList) < len(tasks) {
		t.Errorf("Ожидалось минимум %d задач, получено %d", len(tasks), len(tasksList))
	}

	found := 0
	for _, id := range ids {
		for _, task := range tasksList {
			if task["id"] == id {
				found++
				break
			}
		}
	}

	if found != len(ids) {
		t.Errorf("Найдено только %d из %d созданных задач", found, len(ids))
	}
}
