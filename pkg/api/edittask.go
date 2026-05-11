package api

import (
    "encoding/json"
    "net/http"

    "final_project/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    if id == "" {
        writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
        return
    }

    task, err := db.GetTask(id)
    if err != nil {
        writeJSON(w, map[string]string{"error": err.Error()})
        return
    }

    writeJSON(w, task)
}

// updateTaskHandler обрабатывает PUT /api/task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
    var task db.Task
    if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
        writeJSON(w, map[string]string{"error": "ошибка десериализации JSON"})
        return
    }

    if task.ID == "" {
        writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
        return
    }

    if task.Title == "" {
        writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"})
        return
    }

    if err := checkDate(&task); err != nil {
        writeJSON(w, map[string]string{"error": err.Error()})
        return
    }

    if err := db.UpdateTask(&task); err != nil {
        writeJSON(w, map[string]string{"error": err.Error()})
        return
    }

    writeJSON(w, map[string]struct{}{})
}