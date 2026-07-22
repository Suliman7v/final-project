package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTask(t *testing.T) {
	cleanupDB(t)

	now := time.Now()

	task := task{
		date:    now.Format("20060102"),
		title:   "Созвон в 16:00",
		comment: "Обсуждение планов",
		repeat:  "d 5",
	}

	todo := addTask(t, task)

	body, err := requestJSON("api/task", nil, http.MethodGet)
	assert.NoError(t, err)
	var m map[string]string
	err = json.Unmarshal(body, &m)
	assert.NoError(t, err)

	e, ok := m["error"]
	assert.False(t, !ok || len(fmt.Sprint(e)) == 0,
		"Ожидается ошибка для вызова /api/task")

	body, err = requestJSON("api/task?id="+todo, nil, http.MethodGet)
	assert.NoError(t, err)
	err = json.Unmarshal(body, &m)
	assert.NoError(t, err)

	assert.Equal(t, todo, m["id"])
	assert.Equal(t, task.date, m["date"])
	assert.Equal(t, task.title, m["title"])
	assert.Equal(t, task.comment, m["comment"])
	assert.Equal(t, task.repeat, m["repeat"])
}

type fulltask struct {
	id string
	task
}

func TestEditTask(t *testing.T) {
	cleanupDB(t)

	createTask := map[string]interface{}{
		"date":    time.Now().Format("20060102"),
		"title":   "Заказать хинкали",
		"comment": "в 18:00",
		"repeat":  "d 7",
	}

	id := CreateTask(t, createTask)

	updateTask := map[string]interface{}{
		"id":      id,
		"date":    "20261220",
		"title":   "Заказать осетинские пироги",
		"comment": "с сыром",
		"repeat":  "d 14",
	}

	body, _ := json.Marshal(updateTask)
	req, _ := http.NewRequest("PUT", "http://localhost:7540/api/task", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", resp.StatusCode)
	}

	updated := GetTask(t, id)

	if updated["title"] != "Заказать осетинские пироги" {
		t.Errorf("Title не обновился: %v", updated["title"])
	}
	if updated["comment"] != "с сыром" {
		t.Errorf("Comment не обновился: %v", updated["comment"])
	}
	if updated["repeat"] != "d 14" {
		t.Errorf("Repeat не обновился: %v", updated["repeat"])
	}
}
