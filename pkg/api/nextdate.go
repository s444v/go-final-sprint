package api

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Фукнция для поиска след. даты по заданному правилу
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	var nextDate string
	dstartTime, err := time.Parse(TIMEFORMAT, dstart)
	if err != nil {
		return "", fmt.Errorf("cant parse start_date into time.Time %w", err)
	}
	rule := strings.Split(repeat, " ")
	if err = validateRule(rule); err != nil {
		return "", fmt.Errorf("wrong format: %w", err)
	}
	switch rule[0] {
	case "d":
		nextDate, err = dRule(rule, dstartTime, now)
	case "y":
		nextDate = yRule(dstartTime, now)
	case "w":
		nextDate, err = wRule(rule, dstartTime, now)
	case "m":
		nextDate, err = mRule(rule, dstartTime, now)
	default:
		return "", errors.New("wrong rule")
	}
	if err != nil {
		return "", fmt.Errorf("cant find next date %w", err)
	}
	return nextDate, err
}
