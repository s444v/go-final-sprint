package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse(TIMEFORMAT, r.FormValue("now"))
	if err != nil {
		http.Error(w, "cant parse time", http.StatusBadRequest)
		return
	}
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	result, err := nextDate(now, date, repeat)
	if err != nil {
		http.Error(w, "cant find next date", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

func nextDate(now time.Time, dstart string, repeat string) (string, error) {
	rule := strings.Split(repeat, " ")
	nextDate, err := time.Parse(TIMEFORMAT, dstart)
	if err != nil {
		return "", err
	}
	switch rule[0] {
	case "d":
		if len(rule) != 2 {
			return "", fmt.Errorf("wrong rule format")
		}
		interval, err := strconv.Atoi(rule[1])
		if err != nil {
			return "", fmt.Errorf("conv error %w", err)
		}
		if interval > 400 || interval < 0 {
			return "", fmt.Errorf("wrong number of days %w", err)
		}
		for {
			nextDate = nextDate.AddDate(0, 0, interval)
			if afterNow(nextDate, now) {
				break
			}
		}
	case "y":
		for {
			nextDate = nextDate.AddDate(1, 0, 0)
			if afterNow(nextDate, now) {
				break
			}
		}
	case "w":
		if len(rule) != 2 {
			return "", fmt.Errorf("wrong rule format")
		}
		weekDays := strings.Split(rule[1], ",")
		nowWeekDay := now.Weekday()
		dist := 7
		for _, v := range weekDays {
			weekDay, err := strconv.Atoi(v)
			if err != nil {
				return "", fmt.Errorf("conv error %w", err)
			}
			if weekDay < 1 || weekDay > 7 {
				return "", fmt.Errorf("wrong day of the week")
			}
			weekDay = weekDay % 7
			if dist > (weekDay-int(nowWeekDay)+7)%8 {
				dist = (weekDay - int(nowWeekDay) + 7) % 8
			}
		}
		nextDate = now.AddDate(0, 0, dist)
	case "m":
		if len(rule) != 2 || len(rule) != 3 {
			return "", fmt.Errorf("wrong rule format")
		}
	default:
		return "", errors.New("wrong rule")
	}
	return nextDate.Format(TIMEFORMAT), nil
}

func afterNow(date, now time.Time) bool {
	return date.After(now)
}
