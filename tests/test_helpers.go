package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"final-project/pkg/db"
)

func CreateTask(t *testing.T, task map[string]interface{}) string {
	body, _ := json.Marshal(task)
	resp, err := http.Post("http://localhost:7540/api/task", "application/json", strings.NewReader(string(body)))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался статус 200, получен %d", resp.StatusCode)
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatal(err)
	}

	id, ok := result["id"].(string)
	if !ok {
		t.Fatal("ID не получен")
	}
	return id
}

func GetTask(t *testing.T, id string) map[string]interface{} {
	resp, err := http.Get(fmt.Sprintf("http://localhost:7540/api/task?id=%s", id))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var task map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&task)
	if err != nil {
		t.Fatal(err)
	}
	return task
}

func InitTestDB(t *testing.T) {
	err := db.Init(TestDBConnStr)
	if err != nil {
		t.Fatalf("Не удалось подключиться к тестовой БД: %v", err)
	}
}

func CloseTestDB(t *testing.T) {
	db.Close()
}
