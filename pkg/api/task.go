package api

import (
	"encoding/json"
	"final-project/pkg/db"
	"net/http"
	"time"
)

// Вспомогательная функция для возврата JSON ответа
func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Вспомогательная функция для возврата ошибки в JSON
func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	response := map[string]string{
		"error": message,
	}
	writeJSON(w, response, statusCode)
}

// Вспомогательная функция для сравнения дней
func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleAddTask(w, r)
	case http.MethodGet:
		handleGetTask(w, r)
	case http.MethodPut:
		handleUpdateTask(w, r)
	case http.MethodDelete:
		handleDeleteTask(w, r) // Добавляем DELETE
	default:
		writeJSONError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func handleAddTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Date    string `json:"date"`
		Title   string `json:"title"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		writeJSONError(w, "Ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		writeJSONError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	now := time.Now()

	if req.Date == "" {
		req.Date = now.Format("20060102")
	}

	taskDate, err := time.Parse("20060102", req.Date)
	if err != nil {
		writeJSONError(w, "Неверный формат даты", http.StatusBadRequest)
		return
	}

	if req.Repeat != "" {
		_, err := NextDate(now, req.Date, req.Repeat)
		if err != nil {
			writeJSONError(w, "Некорректное правило повторения", http.StatusBadRequest)
			return
		}
		if taskDate.Before(now) && !isSameDay(taskDate, now) {
			next, err := NextDate(now, req.Date, req.Repeat)
			if err != nil {
				writeJSONError(w, "Некорректное правило повторения", http.StatusBadRequest)
				return
			}
			req.Date = next
		}
	} else if taskDate.Before(now) && !isSameDay(taskDate, now) {
		req.Date = now.Format("20060102")
	}

	task := &db.Task{
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	id, err := db.AddTask(task)
	if err != nil {
		writeJSONError(w, "Ошибка при добавлении задачи в БД", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id": id,
	}
	writeJSON(w, response, http.StatusOK)
}

func handleGetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSONError(w, "Ошибка при получении задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if task == nil {
		writeJSONError(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	writeJSON(w, task, http.StatusOK)
}

func handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID      string `json:"id"`
		Date    string `json:"date"`
		Title   string `json:"title"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		writeJSONError(w, "Ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	if req.ID == "" {
		writeJSONError(w, "Не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		writeJSONError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	existingTask, err := db.GetTask(req.ID)
	if err != nil {
		writeJSONError(w, "Ошибка при проверке задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if existingTask == nil {
		writeJSONError(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	now := time.Now()

	if req.Date == "" {
		req.Date = now.Format("20060102")
	}

	taskDate, err := time.Parse("20060102", req.Date)
	if err != nil {
		writeJSONError(w, "Неверный формат даты", http.StatusBadRequest)
		return
	}

	if req.Repeat != "" {
		_, err := NextDate(now, req.Date, req.Repeat)
		if err != nil {
			writeJSONError(w, "Некорректное правило повторения", http.StatusBadRequest)
			return
		}

		if taskDate.Before(now) && !isSameDay(taskDate, now) {
			next, err := NextDate(now, req.Date, req.Repeat)
			if err != nil {
				writeJSONError(w, "Некорректное правило повторения", http.StatusBadRequest)
				return
			}
			req.Date = next
		}
	} else if taskDate.Before(now) && !isSameDay(taskDate, now) {
		req.Date = now.Format("20060102")
	}

	task := &db.Task{
		ID:      req.ID,
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}
	err = db.UpdateTask(task)
	if err != nil {
		writeJSONError(w, "Ошибка при обновлении задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}

func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeJSONError(w, "Ошибка при удалении задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}
