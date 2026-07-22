package api

import (
	"net/http"

	"final-project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

const limit = 50

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(limit, "")
	if err != nil {
		writeJSONError(w, "Ошибка при получении задач", http.StatusInternalServerError)
		return
	}

	resp := TasksResp{
		Tasks: tasks,
	}

	writeJSON(w, resp, http.StatusOK)
}
