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

func getTasks(t *testing.T, search string) []map[string]string {
	url := "api/tasks"
	if Search {
		url += "?search=" + search
	}
	body, err := requestJSON(url, nil, http.MethodGet)
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

	// Создаем задачи и сохраняем их ID
	var ids []string
	for _, task := range tasks {
		id := CreateTask(t, task)
		ids = append(ids, id)
	}

	// Получаем список всех задач
	resp, err := http.Get("http://localhost:7540/api/tasks")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", resp.StatusCode)
	}

	// Проверяем структуру ответа
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatal(err)
	}

	// Проверяем, что есть поле tasks
	tasksResp, ok := result["tasks"].([]interface{})
	if !ok {
		t.Fatal("Нет поля tasks или оно не массив")
	}

	// Проверяем, что количество задач не меньше созданных
	if len(tasksResp) < len(tasks) {
		t.Errorf("Ожидалось минимум %d задач, получено %d", len(tasks), len(tasksResp))
	}

	// Проверяем, что созданные задачи есть в списке
	found := 0
	for _, id := range ids {
		for _, task := range tasksResp {
			taskMap := task.(map[string]interface{})
			if taskMap["id"] == id {
				found++
				break
			}
		}
	}

	if found != len(ids) {
		t.Errorf("Найдено только %d из %d созданных задач", found, len(ids))
	}
}
