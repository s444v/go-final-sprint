package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/s444v/go-final-sprint/pkg/database"
)

// Функция для записи в responseWriter v interface в json
func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(v)
}

// Функция для проверки даты
func checkDate(task *database.Task) error {
	now := time.Now()
	// если дата отсутствует, то берем сегодняшнюю дату
	if task.Date == "" {
		task.Date = now.Format(TIMEFORMAT)
	}
	t, err := time.Parse(TIMEFORMAT, task.Date)
	if err != nil {
		return fmt.Errorf("дата представлена в формате, отличном от 20060102, %w", err)
	}
	// проверяем чтоб дата была после now
	if after(now, t) {
		if len(task.Repeat) == 0 {
			// если правила повторения нет, то берём сегодняшнее число
			task.Date = now.Format(TIMEFORMAT)
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			// в противном случае, берём вычисленную ранее следующую дату
			task.Date = next
		}
	}
	return err
}

/*
Функция для проверки даты
функция проверяет можно ли создать дату с определенным днем
например 2025.02.30 = false
например 2025.02.20 = true
*/
func isValidDate(baseDate time.Time, day int) bool {
	year := baseDate.Year()
	month := baseDate.Month()
	if day < 0 {
		day = daysInMonth(year, month) + day + 1
	}
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return t.Year() == year && int(t.Month()) == int(month) && t.Day() == day
}

// Функция возвращает количество дней в месяце
func daysInMonth(year int, month time.Month) int {
	t := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	t = t.AddDate(0, 1, -1)
	return t.Day()
}

// проверка на то что data1 после data2
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

// проверка на то что даты равны
func equal(date1, date2 time.Time) bool {
	dy, dm, dd := date1.Date()
	ny, nm, nd := date2.Date()

	return dy == ny && dm == nm && dd == nd
}

// проверка формата правил
func validateRule(rule []string) error {
	if len(rule) == 0 {
		return errors.New("пустой массив rule")
	}
	letter := rule[0]
	switch letter {
	case "d", "w":
		if len(rule) != 2 {
			return fmt.Errorf("для '%s' ожидается 2 элемента, а получено %d", letter, len(rule))
		}
	case "y":
		if len(rule) != 1 {
			return fmt.Errorf("для 'y' ожидается 1 элемент, а получено %d", len(rule))
		}
	case "m":
		if len(rule) != 2 && len(rule) != 3 {
			return fmt.Errorf("для 'm' ожидается 2 или 3 элемента, а получено %d", len(rule))
		}
	default:
		return fmt.Errorf("недопустимое значение в rule: '%s'", letter)
	}

	return nil
}
