package api

import (
    "net/http"
    "time"

    "final_project/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        return
    }

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

    if task.Repeat == "" {
        if err := db.DeleteTask(id); err != nil {
            writeJSON(w, map[string]string{"error": err.Error()})
            return
        }
    } else {
        now := time.Now()
        nextDate, err := NextDate(now, task.Date, task.Repeat)
        if err != nil {
            writeJSON(w, map[string]string{"error": err.Error()})
            return
        }
        if err := db.UpdateTaskDate(id, nextDate); err != nil {
            writeJSON(w, map[string]string{"error": err.Error()})
            return
        }
    }

    writeJSON(w, map[string]struct{}{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        return
    }

    id := r.URL.Query().Get("id")
    if id == "" {
        writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
        return
    }

    if err := db.DeleteTask(id); err != nil {
        writeJSON(w, map[string]string{"error": err.Error()})
        return
    }

    writeJSON(w, map[string]struct{}{})
}