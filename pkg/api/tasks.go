package api

import (
	"final-project/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50, "")
	if err != nil {
		writeJSONError(w, "Ошибка при получении задач", http.StatusInternalServerError)
		return
	}

	resp := TasksResp{
		Tasks: tasks,
	}

	writeJSON(w, resp, http.StatusOK)
}
