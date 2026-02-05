package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func afterNow(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()
	return y1 > y2 || (y1 == y2 && m1 > m2) || (y1 == y2 && m1 == m2 && d1 > d2)
}

func weekdayToNumber(wd time.Weekday) int {
	if wd == time.Sunday {
		return 7
	}
	return int(wd)
}

func nextdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeatStr := r.URL.Query().Get("repeat")

	if nowStr == "" || dateStr == "" || repeatStr == "" {
		http.Error(w, "Отсутствует обязательный параметр", http.StatusBadRequest)
		return
	}

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "неверный формат даты", http.StatusBadRequest)
		return
	}

	result, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("параметр repeat не может быть пустым")
	}

	parsedDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("неверный формат даты")
	}

	parts := strings.Split(repeat, " ")
	if len(parts) <= 0 {
		return "", errors.New("неверный формат правила")
	}

	ruleType := parts[0]
	switch ruleType {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("неверный формат для правила 'd'")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("неверный интервал для правила 'd'")
		}
		if interval < 1 || interval > 400 {
			return "", errors.New("интервал должен быть от 1 до 400")
		}

		next := parsedDate
		for {
			next = next.AddDate(0, 0, interval)
			if afterNow(next, now) {
				break
			}
		}
		return next.Format("20060102"), nil

	case "y":
		if len(parts) != 1 {
			return "", errors.New("неверный формат для правила 'y'")
		}
		next := parsedDate
		for {
			next = next.AddDate(1, 0, 0)
			if afterNow(next, now) {
				break
			}
		}
		return next.Format("20060102"), nil

	case "w":
		if len(parts) != 2 {
			return "", errors.New("неверный формат для правила 'w'")
		}

		dayStrs := strings.Split(parts[1], ",")
		allowedDays := make(map[int]bool)
		for _, dayStr := range dayStrs {
			day, err := strconv.Atoi(strings.TrimSpace(dayStr))
			if err != nil || day < 1 || day > 7 {
				return "", errors.New("недопустимый день недели")
			}
			allowedDays[day] = true
		}

		next := parsedDate
		for {
			next = next.AddDate(0, 0, 1)
			weekdayNum := weekdayToNumber(next.Weekday())
			if allowedDays[weekdayNum] && afterNow(next, now) {
				break
			}
		}
		return next.Format("20060102"), nil

	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", errors.New("неверный формат для правила 'm'")
		}

		dayStrs := strings.Split(parts[1], ",")
		allowedDays := make(map[int]bool)
		for _, dayStr := range dayStrs {
			dayM, err := strconv.Atoi(strings.TrimSpace(dayStr))
			if err != nil {
				return "", errors.New("недопустимый день")
			}
			if (dayM < -2 || dayM > 31) || dayM == 0 {
				return "", errors.New("недопустимый день")
			}
			allowedDays[dayM] = true
		}

		if len(parts) == 3 {
			monthStrs := strings.Split(parts[2], ",")
			for _, monthStr := range monthStrs {
				month, err := strconv.Atoi(strings.TrimSpace(monthStr))
				if err != nil || month < 1 || month > 12 {
					return "", errors.New("недопустимый месяц")
				}
			}
		}
		// Пока не реализовано
		return "", errors.New("правило 'm' пока не реализовано")

	default:
		return "", errors.New("неподдерживаемый формат правила")
	}
}
