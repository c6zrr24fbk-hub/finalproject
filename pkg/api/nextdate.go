package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("правило повторения не указано")
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", err
	}

	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	parts := strings.Split(repeat, " ")
	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("неверный формат: d <число>")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("интервал должен быть числом от 1 до 400")
		}
		for {
			date = date.AddDate(0, 0, days)
			if date.After(now) {
				break
			}
		}
		return date.Format(DateFormat), nil

	case "y":
		if len(parts) != 1 {
			return "", errors.New("правило y не требует параметров")
		}
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				break
			}
		}
		return date.Format(DateFormat), nil

	case "w":
		if len(parts) != 2 {
			return "", errors.New("неверный формат: w <дни недели>")
		}
		wdays := strings.Split(parts[1], ",")
		var weekdays [7]bool
		for _, s := range wdays {
			d, err := strconv.Atoi(s)
			if err != nil || d < 1 || d > 7 {
				return "", errors.New("день недели должен быть от 1 до 7")
			}
			idx := d % 7
			weekdays[idx] = true
		}
		for {
			if weekdays[date.Weekday()] && date.After(now) {
				break
			}
			date = date.AddDate(0, 0, 1)
		}
		return date.Format(DateFormat), nil

	case "m":
		type monthRule struct {
			days   [32]bool
			last    bool
			prelast bool
			months  [13]bool
		}
		rule := monthRule{}
		if len(parts) < 2 || len(parts) > 3 {
			return "", errors.New("неверный формат: m <дни> [месяцы]")
		}
		dayStrs := strings.Split(parts[1], ",")
		for _, s := range dayStrs {
			switch s {
			case "-1":
				rule.last = true
			case "-2":
				rule.prelast = true
			default:
				d, err := strconv.Atoi(s)
				if err != nil || d < 1 || d > 31 {
					return "", errors.New("неверный день месяца")
				}
				rule.days[d] = true
			}
		}
		if len(parts) == 3 {
			monthStrs := strings.Split(parts[2], ",")
			for _, s := range monthStrs {
				m, err := strconv.Atoi(s)
				if err != nil || m < 1 || m > 12 {
					return "", errors.New("неверный месяц")
				}
				rule.months[m] = true
			}
		} else {
			for i := 1; i <= 12; i++ {
				rule.months[i] = true
			}
		}
		for {
			if !rule.months[int(date.Month())] {
				date = date.AddDate(0, 0, 1)
				continue
			}
			day := date.Day()
			lastDay := lastDayOfMonth(date)

			matched := false
			if rule.last && day == lastDay {
				matched = true
			} else if rule.prelast && day == lastDay-1 {
				matched = true
			} else if rule.days[day] {
				matched = true
			}
			if matched && date.After(now) {
				break
			}
			date = date.AddDate(0, 0, 1)
		}
		return date.Format(DateFormat), nil

	default:
		return "", errors.New("неподдерживаемый формат правила")
	}
}

func lastDayOfMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	if dateStr == "" || repeat == "" {
		http.Error(w, "Не указаны параметры date или repeat", http.StatusBadRequest)
		return
	}

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "Неверный формат now", http.StatusBadRequest)
			return
		}
	}

	result, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(result))
}