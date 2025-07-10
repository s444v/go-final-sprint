package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/s444v/go-final-sprint/pkg/database"
)

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(v)
}

func checkDate(task *database.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(TIMEFORMAT)
	}
	t, err := time.Parse(TIMEFORMAT, task.Date)
	if err != nil {
		return fmt.Errorf("дата представлена в формате, отличном от 20060102, %w", err)
	}
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return err
	}
	if after(now, t) {
		if len(task.Repeat) == 0 {
			// если правила повторения нет, то берём сегодняшнее число
			task.Date = now.Format(TIMEFORMAT)
		} else {
			// в противном случае, берём вычисленную ранее следующую дату
			task.Date = next
		}
	}
	return err
}

func findNextValidDate(baseDate time.Time, targetDay int) time.Time {
	year := baseDate.Year()
	month := baseDate.Month()
	check := true
	if targetDay < 0 {
		check = false
	}
	if targetDay <= baseDate.Day() && check {
		month += 1
	} else if (daysInMonth(year, month)+targetDay+1 <= baseDate.Day() || daysInMonth(year, month)+targetDay+1 > daysInMonth(year, month)) && !check {
		month += 1
	}
	for {
		if month > 12 {
			month = 1
			year++
		}
		if targetDay <= daysInMonth(year, month) && check {
			return time.Date(year, month, targetDay, 0, 0, 0, 0, baseDate.Location())
		} else if daysInMonth(year, month)+targetDay+1 <= daysInMonth(year, month) {
			return time.Date(year, month, daysInMonth(year, month)+targetDay+1, 0, 0, 0, 0, baseDate.Location())
		}
		month += 1
	}
}

func daysInMonth(year int, month time.Month) int {
	t := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	t = t.AddDate(0, 1, -1)
	return t.Day()
}

func after(date1, date2 time.Time) bool {
	dy, dm, dd := date1.Date()
	ny, nm, nd := date2.Date()

	if dy != ny {
		return dy > ny
	}

	if dm != nm {
		return dm > nm
	}

	return dd > nd
}
