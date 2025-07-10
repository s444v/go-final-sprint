package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

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
		if after(nextDate, now) {
			break
		}
	}
	return nextDate.Format(TIMEFORMAT), err
}

func wRule(rule []string, dstart string, now time.Time) (string, error) {
	nextDate, err := time.Parse(TIMEFORMAT, dstart)
	if err != nil {
		return "", fmt.Errorf("cant find next date %w", err)
	}
	if !after(nextDate, now) {
		nextDate = now
	}
	if len(rule) != 2 {
		return "", errors.New("wrong rule format")
	}
	weekDays := strings.Split(rule[1], ",")
	nowWeekDay := nextDate.Weekday()
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
	nextDate = nextDate.AddDate(0, 0, dist)
	return nextDate.Format(TIMEFORMAT), nil
}

func mRule(rule []string, dstart string, now time.Time) (string, error) {
	startDate, err := time.Parse(TIMEFORMAT, dstart)
	result := ""
	if err != nil {
		return result, err
	}
	if after(now, startDate) {
		startDate = now
	}
	var nearestDate time.Time
	days := strings.Split(rule[1], ",")
	if len(rule) == 2 {
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
			if after(nearestDate, findNextValidDate(startDate, day)) {
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

func yRule(dstart string, now time.Time) (string, error) {
	nextDate, err := time.Parse(TIMEFORMAT, dstart)
	if err != nil {
		return "", fmt.Errorf("cant find next date %w", err)
	}
	for {
		nextDate = nextDate.AddDate(1, 0, 0)
		if after(nextDate, now) {
			break
		}
	}
	return nextDate.Format(TIMEFORMAT), err
}
