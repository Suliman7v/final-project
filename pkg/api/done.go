package api

import (
	"final-project/pkg/db"
	"net/http"
	"time"
)

func doneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

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

	now := time.Now()

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJSONError(w, "Ошибка при удалении задачи: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSONError(w, "Ошибка при вычислении следующей даты: "+err.Error(), http.StatusBadRequest)
			return
		}

		err = db.UpdateDate(id, nextDate)
		if err != nil {
			writeJSONError(w, "Ошибка при обновлении даты задачи: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}
