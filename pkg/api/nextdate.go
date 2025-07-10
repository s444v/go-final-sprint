package api

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	rule := strings.Split(repeat, " ")
	if rule[0] == "" {
		return "", nil
	}
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
