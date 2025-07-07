package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func nextDayHandler(http.ResponseWriter, *http.Request) {

}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	rule := strings.Split(repeat, " ")
	nextDate, err := time.Parse(timeFormat, dstart)
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
			if weekDay == 7 {
				weekDay--
			}
			if dist > (weekDay-int(nowWeekDay)+7)%7 {
				dist = (weekDay - int(nowWeekDay) + 7) % 7
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
	return nextDate.Format(timeFormat), nil
}

func afterNow(date, now time.Time) bool {
	return date.After(now)
}
