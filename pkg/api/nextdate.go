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
	switch rule[0] {
	case "d":
		nextDate, err := dRule(rule, dstart, now)
		if err != nil {
			return "", fmt.Errorf("cant find next date %w", err)
		}
		return nextDate, err
	case "y":
		nextDate, err := yRule(dstart, now)
		if err != nil {
			return "", fmt.Errorf("cant find next date %w", err)
		}
		return nextDate, err
	case "w":
		nextDate, err := wRule(rule, dstart, now)
		if err != nil {
			return "", fmt.Errorf("cant find next date %w", err)
		}
		return nextDate, err
	case "m":
		nextDate, err := mRule(rule, dstart, now)
		if err != nil {
			return "", fmt.Errorf("cant find next date %w", err)
		}
		return nextDate, err
	default:
		return "", errors.New("wrong rule")
	}
}

func yRule(dstart string, now time.Time) (string, error) {
	nextDate, err := time.Parse(TIMEFORMAT, dstart)
	if err != nil {
		return "", fmt.Errorf("cant find next date %w", err)
	}
	for {
		nextDate = nextDate.AddDate(1, 0, 0)
		if afterNow(nextDate, now) {
			break
		}
	}
	return nextDate.Format(TIMEFORMAT), err
}

func wRule(rule []string, dstart string, now time.Time) (string, error) {
	if len(rule) != 2 {
		return "", errors.New("wrong rule format")
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
			return "", errors.New("wrong day of the week")
		}
		weekDay = weekDay % 7
		if dist > (weekDay-int(nowWeekDay)+7)%8 {
			dist = (weekDay - int(nowWeekDay) + 7) % 8
		}
	}
	nextDate := now.AddDate(0, 0, dist)
	return nextDate.Format(TIMEFORMAT), nil
}

func dRule(rule []string, dstart string, now time.Time) (string, error) {
	nextDate, err := time.Parse(TIMEFORMAT, dstart)
	if err != nil {
		return "", fmt.Errorf("conv error %w", err)
	}
	if len(rule) != 2 {
		return "", errors.New("wrong rule format")
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
	return nextDate.Format(TIMEFORMAT), err
}

func mRule(rule []string, dstart string, now time.Time) (string, error) {
	startDate, err := time.Parse(TIMEFORMAT, dstart)
	result := ""
	if err != nil {
		return result, err
	}
	if afterNow(now, startDate) {
		startDate = now
	}
	var nearestDate time.Time
	days := strings.Split(rule[1], ",")
	if len(rule) == 2 {
		//days := strings.Split(rule[1], ",")
		first := true
		for _, d := range days {
			day, err := strconv.Atoi(d)
			if err != nil {
				return result, fmt.Errorf("conv error %w", err)
			}
			if day > 31 || day < -2 {
				return result, errors.New("wrong day of the month")
			}
			if first {
				nearestDate = findNextValidDate(startDate, day)
				first = false
				continue
			}
			if afterNow(nearestDate, findNextValidDate(startDate, day)) {
				nearestDate = findNextValidDate(startDate, day)
			}
		}
		result = nearestDate.Format(TIMEFORMAT)
	} else if len(rule) == 3 {
		months := strings.Split(rule[2], ",")
		found := false
		for _, m := range months {
			month, err := strconv.Atoi(m)
			if err != nil {
				return result, fmt.Errorf("conv error %w", err)
			}
			if month > 12 || month < 1 {
				return result, errors.New("wrong month")
			}
			for _, d := range days {
				day, err := strconv.Atoi(d)
				if err != nil {
					return result, fmt.Errorf("conv error %w", err)
				}
				if day > 31 || day < -2 {
					return result, errors.New("wrong day of the month")
				}
				d := time.Date(startDate.Year(), time.Month(month), day, 0, 0, 0, 0, time.UTC)
				if d.Month() != time.Month(month) || d.Day() != day {
					continue
				}
				if (d.Equal(startDate) || d.After(startDate)) && (!found || d.Before(nearestDate)) {
					nearestDate = d
					found = true
				}
			}
			result = nearestDate.Format(TIMEFORMAT)
		}
	} else {
		return result, errors.New("wrong rule format")
	}
	return result, nil
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
	// Берем 1-е число следующего месяца, вычитаем 1 день
	t := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	t = t.AddDate(0, 1, 0)
	t = t.AddDate(0, 0, -1)
	return t.Day()
}

func afterNow(date, now time.Time) bool {
	return date.After(now)
}
