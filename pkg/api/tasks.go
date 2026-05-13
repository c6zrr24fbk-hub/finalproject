package api

import (
	"net/http"

	"final_project/pkg/db"
)

const tasksLimit = 50

type TasksResponse struct {
	Tasks []db.Task `json:"tasks"`
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")

	tasks, err := db.GetTasks(tasksLimit, search)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, TasksResponse{Tasks: tasks})
}