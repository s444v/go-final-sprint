package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func dRule(rule []string, dstart, now time.Time) (string, error) {
	nextDate := dstart
	interval, err := strconv.Atoi(rule[1])
	if err != nil {
		return "", fmt.Errorf("cant convert interval to int %w", err)
	}
	if interval > 400 || interval < 1 {
		return "", fmt.Errorf("interval must be <= 400 or >= 1 %w", err)
	}
	for {
		nextDate = nextDate.AddDate(0, 0, interval)
		if after(nextDate, now) {
			break
		}
	}
	return nextDate.Format(TIMEFORMAT), err
}

func wRule(rule []string, dstart, now time.Time) (string, error) {
	nextDate := dstart
	if !after(nextDate, now) {
		nextDate = now
	}
	weekDays := strings.Split(rule[1], ",")
	nowWeekDay := int(nextDate.Weekday())
	interval := 7
	for _, v := range weekDays {
		weekDay, err := strconv.Atoi(v)
		if err != nil {
			return "", fmt.Errorf("cant convert weekday to int %w", err)
		}
		if weekDay < 1 || weekDay > 7 {
			return "", errors.New("weekday must be <= 7 or >= 1")
		}
		if interval > (weekDay-nowWeekDay+6)%7 {
			interval = (weekDay - nowWeekDay + 6) % 7
		}
	}
	nextDate = nextDate.AddDate(0, 0, interval+1)
	return nextDate.Format(TIMEFORMAT), nil
}

func mRule(rule []string, dstart, now time.Time) (string, error) {
	if after(dstart, now) {
		now = dstart
	}
	if len(rule) == 2 {
		nearestDate := now
		first := true
		for _, d := range strings.Split(rule[1], ",") {
			day, err := strconv.Atoi(d)
			if err != nil {
				return "", fmt.Errorf("cant convert day to int %w", err)
			}
			tmp := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
			for {
				if isValidDate(tmp, day) {
					if day < 0 {
						day = daysInMonth(tmp.Year(), tmp.Month()) + day + 1
					}
					tmp = time.Date(tmp.Year(), tmp.Month(), day, 0, 0, 0, 0, time.UTC)
					if after(tmp, now) {
						break
					}
				}
				tmp = time.Date(tmp.Year(), tmp.Month()+1, 1, 0, 0, 0, 0, time.UTC)
			}
			fmt.Println(nearestDate, tmp)
			if first {
				nearestDate = tmp
				first = false
				continue
			}
			if after(nearestDate, tmp) {
				nearestDate = tmp
			}
		}
		return nearestDate.Format(TIMEFORMAT), nil
	}
	return "", nil
}

// func mRule(rule []string, dstart, now time.Time) (string, error) {
// 	startDate := dstart
// 	if after(now, startDate) {
// 		startDate = now
// 	}
// 	var nearestDate time.Time
// 	if len(rule) == 2 {
// 		nearestDate, err := find(strings.Split(rule[1], ","), startDate)
// 		return nearestDate.Format(TIMEFORMAT), err
// 	}
// 	months := strings.Split(rule[2], ",")
// 	for _, m := range months {
// 		month, err := strconv.Atoi(m)
// 		if err != nil {
// 			return nearestDate.Format(TIMEFORMAT), fmt.Errorf("cant convert month to int %w", err)
// 		}
// 		if month > 12 || month < 1 {
// 			return nearestDate.Format(TIMEFORMAT), errors.New("month must be <= 12 or >= 1")
// 		}
// 		//tmp := time.Date(startDate.Year(), time.Month(month), 1, 0, 0, 0, 0, time.UTC)
// 		for _, d := range strings.Split(rule[1], ",") {
// 			day, err := strconv.Atoi(d)
// 			if err != nil {
// 				return nearestDate.Format(TIMEFORMAT), fmt.Errorf("cant convert day to int %w", err)
// 			}
// 			if day > 31 || day < -2 {
// 				return nearestDate.Format(TIMEFORMAT), errors.New("day of the month must be <= 31 or >= -2")
// 			}
// 		}
// 		nearestDate, err = find(strings.Split(rule[1], ","), startDate)
// 		if err != nil {
// 			return "", fmt.Errorf("cant convert month to int %w", err) // изменить текст ошибки
// 		}
// 	}
// 	return nearestDate.Format(TIMEFORMAT), nil
// }

func yRule(dstart, now time.Time) string {
	nextDate := dstart
	for {
		nextDate = nextDate.AddDate(1, 0, 0)
		if after(nextDate, now) {
			break
		}
	}
	return nextDate.Format(TIMEFORMAT)
}

// func find(days []string, startDate time.Time) (time.Time, error) {
// 	var nearestDate time.Time
// 	first := true
// 	for _, d := range days {
// 		day, err := strconv.Atoi(d)
// 		if err != nil {
// 			return nearestDate, fmt.Errorf("cant convert day to int %w", err)
// 		}
// 		if day > 31 || day < -2 {
// 			return nearestDate, errors.New("day of the month must be <= 31 or >= -2")
// 		}
// 		tmp := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.UTC)
// 		for {
// 			if !isValidDate(tmp, day) {
// 				tmp = time.Date(tmp.Year(), tmp.Month()+1, 1, 0, 0, 0, 0, time.UTC)
// 				continue
// 			}
// 			if day < 0 {
// 				day = daysInMonth(tmp.Year(), tmp.Month()) + day + 1
// 			}
// 			tmp = time.Date(tmp.Year(), tmp.Month(), day, 0, 0, 0, 0, time.UTC)
// 			if after(tmp, startDate) {
// 				break
// 			}
// 			tmp = time.Date(tmp.Year(), tmp.Month()+1, 1, 0, 0, 0, 0, time.UTC)
// 		}
// 		if first {
// 			nearestDate = tmp
// 			first = false
// 			continue
// 		}
// 		if after(nearestDate, tmp) {
// 			nearestDate = tmp
// 		}
// 	}
// 	return nearestDate, nil
// }

// func find(year, month int, targetDay string) bool {
// 		day, err := strconv.Atoi(targetDay)
// 		if err != nil {
// 			return nearestDate, fmt.Errorf("cant convert day to int %w", err)
// 		}
// 		if day > 31 || day < -2 {
// 			return nearestDate, errors.New("day of the month must be <= 31 or >= -2")
// 		}
// 		tmp := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
// 		return isValidDate(tmp, day)
// 		for {
// 			if !isValidDate(tmp, day) {
// 				tmp = time.Date(tmp.Year(), tmp.Month()+1, 1, 0, 0, 0, 0, time.UTC)
// 				continue
// 			}
// 			if day < 0 {
// 				day = daysInMonth(tmp.Year(), tmp.Month()) + day + 1
// 			}
// 			tmp = time.Date(tmp.Year(), tmp.Month(), day, 0, 0, 0, 0, time.UTC)
// 			if after(tmp, startDate) {
// 				break
// 			}
// 			tmp = time.Date(tmp.Year(), tmp.Month()+1, 1, 0, 0, 0, 0, time.UTC)
// 		}
// 		if first {
// 			nearestDate = tmp
// 			first = false
// 			continue
// 		}
// 		if after(nearestDate, tmp) {
// 			nearestDate = tmp
// 		}
// 	}
// 	return nearestDate, nil
// }
