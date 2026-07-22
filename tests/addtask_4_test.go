package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"testing"
)

func requestJSON(apipath string, values map[string]any, method string) ([]byte, error) {
	var (
		data []byte
		err  error
	)

	if len(values) > 0 {
		data, err = json.Marshal(values)
		if err != nil {
			return nil, err
		}
	}
	var resp *http.Response

	req, err := http.NewRequest(method, getURL(apipath), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	if len(Token) > 0 {
		jar, err := cookiejar.New(nil)
		if err != nil {
			return nil, err
		}
		jar.SetCookies(req.URL, []*http.Cookie{
			{
				Name:  "token",
				Value: Token,
			},
		})
		client.Jar = jar
	}

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}
	return io.ReadAll(resp.Body)
}

func postJSON(apipath string, values map[string]any, method string) (map[string]any, error) {
	var (
		m   map[string]any
		err error
	)

	body, err := requestJSON(apipath, values, method)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &m)
	return m, err
}

type task struct {
	date    string
	title   string
	comment string
	repeat  string
}

func TestAddTask(t *testing.T) {
	cleanupDB(t)

	tests := []struct {
		name    string
		task    map[string]interface{}
		wantErr bool
	}{
		{
			name: "Пустой заголовок",
			task: map[string]interface{}{
				"title": "",
			},
			wantErr: true,
		},
		{
			name: "Нормальная задача",
			task: map[string]interface{}{
				"date":    "20261220",
				"title":   "Сделать что-нибудь",
				"comment": "Хорошо отдохнуть",
			},
			wantErr: false,
		},
		{
			name: "Задача с повторением",
			task: map[string]interface{}{
				"date":    "20260228",
				"title":   "Фитнес",
				"comment": "",
				"repeat":  "d 1",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Пытаемся создать задачу
			body, _ := json.Marshal(tt.task)
			resp, err := http.Post("http://localhost:7540/api/task", "application/json", bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if tt.wantErr {
				// Ожидаем ошибку
				if resp.StatusCode != http.StatusBadRequest {
					t.Errorf("Ожидался статус 400, получен %d", resp.StatusCode)
				}
			} else {
				// Ожидаем успех
				if resp.StatusCode != http.StatusOK {
					t.Errorf("Ожидался статус 200, получен %d", resp.StatusCode)
				}

				// Проверяем, что ID вернулся
				var result map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&result)
				if err != nil {
					t.Fatal(err)
				}

				id, ok := result["id"].(string)
				if !ok || id == "" {
					t.Error("ID не получен или пустой")
				}
			}
		})
	}
}
