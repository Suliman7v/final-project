package tests

import (
	"encoding/json"
	"net/http"
	"testing"
)

func cleanupDB(t *testing.T) {
	resp, err := http.Get("http://localhost:7540/api/tasks")
	if err != nil {
		t.Fatal("Не удалось получить список задач:", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatal("Не удалось декодировать ответ:", err)
	}

	tasks, ok := result["tasks"].([]interface{})
	if !ok {
		return
	}

	for _, task := range tasks {
		taskMap, ok := task.(map[string]interface{})
		if !ok {
			continue
		}

		id, ok := taskMap["id"].(string)
		if !ok {
			continue
		}

		req, err := http.NewRequest("DELETE", "http://localhost:7540/api/task?id="+id, nil)
		if err != nil {
			t.Logf("Ошибка при создании запроса для задачи %s: %v", id, err)
			continue
		}

		client := &http.Client{}
		delResp, err := client.Do(req)
		if err != nil {
			t.Logf("Ошибка при удалении задачи %s: %v", id, err)
			continue
		}
		delResp.Body.Close()
	}
}
