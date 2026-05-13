package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"final_project/pkg/db"
)

func afterNow(date, now time.Time) bool {
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return date.After(now)
}

func checkDate(task *db.Task) error {
	now := time.Now()
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if task.Date == "" {
		task.Date = nowDate.Format(DateFormat)
	}

	parsedDate, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return errors.New("неверный формат даты")
	}

	if task.Repeat != "" {
		next, err := NextDate(nowDate, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		if parsedDate.Before(nowDate) || parsedDate.Equal(nowDate) {
			task.Date = next
		}
	} else {
		if parsedDate.Before(nowDate) || parsedDate.Equal(nowDate) {
			task.Date = nowDate.Format(DateFormat)
		}
	}
	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ошибка десериализации JSON"})
		return
	}

	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"id": idToString(id)})
}

func idToString(id int64) string {
	return strconv.FormatInt(id, 10)
}