package api

import (
	"net/http"

	"final_project/pkg/db"
)

type TasksResponse struct {
	Tasks []db.Task `json:"tasks"`
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")
	limit := 50 

	tasks, err := db.GetTasks(limit, search)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, TasksResponse{Tasks: tasks})
}