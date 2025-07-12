package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

/*
Функция для поиска след. даты по правилу d (интервал дней)
*/
func dRule(rule []string, dstart, now time.Time) (string, error) {
	nextDate := dstart
	interval, err := strconv.Atoi(rule[1])
	if err != nil {
		return "", fmt.Errorf("cant convert interval to int %w", err)
	}
	// интервал может быть от 1 до 400
	if interval > 400 || interval < 1 {
		return "", fmt.Errorf("interval must be <= 400 or >= 1 %w", err)
	}
	for {
		//добавляется интервал
		nextDate = nextDate.AddDate(0, 0, interval)
		//пока дата не больше сейчас
		if after(nextDate, now) {
			break
		}
	}
	return nextDate.Format(TIMEFORMAT), err
}

/*
Функция для поиска след. даты по правилу w (по дням неделям)
*/
func wRule(rule []string, dstart, now time.Time) (string, error) {
	nextDate := dstart
	if !after(nextDate, now) {
		nextDate = now
	}
	weekDays := strings.Split(rule[1], ",")
	nowWeekDay := int(nextDate.Weekday())
	interval := 7
	//идет цикл по дням недели
	for _, v := range weekDays {
		weekDay, err := strconv.Atoi(v)
		if err != nil {
			return "", fmt.Errorf("cant convert weekday to int %w", err)
		}
		if weekDay < 1 || weekDay > 7 {
			return "", errors.New("weekday must be <= 7 or >= 1")
		}
		//находится самый короткий интервал от сегодняшнего дня недели
		if interval > (weekDay-nowWeekDay+6)%7 {
			interval = (weekDay - nowWeekDay + 6) % 7
		}
	}
	//к сегодня добавляется интервал
	nextDate = nextDate.AddDate(0, 0, interval+1)
	return nextDate.Format(TIMEFORMAT), nil
}

/*
Функция для поиска след. даты по правилу m (по дням месяца)
*/
func mRule(rule []string, dstart, now time.Time) (string, error) {
	// если dstart> now , то now = dstart
	// сделано если человек выполнил задачу раньше срока
	if after(dstart, now) {
		now = dstart
	}
	// если месяцы не указаны
	if len(rule) == 2 {
		nearestDate := now
		first := true
		// цикл по дням
		for _, d := range strings.Split(rule[1], ",") {
			day, err := strconv.Atoi(d)
			if err != nil {
				return "", fmt.Errorf("cant convert day to int %w", err)
			}
			if day > 31 || day < -2 {
				return "", errors.New("invalid day number")
			}
			// создаем дату с сегодняшним годом и месяцем
			tmp := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
			// начинаем цикл пока tmp < now
			for {
				// проверяем можно ли создать такую дату
				if isValidDate(tmp, day) {
					//если -2, то берем предпоследний день
					// если -1, то берем последний
					if day < 0 {
						day = daysInMonth(tmp.Year(), tmp.Month()) + day + 1
					}
					// создаем дату и проверяем
					tmp = time.Date(tmp.Year(), tmp.Month(), day, 0, 0, 0, 0, time.UTC)
					if after(tmp, now) {
						break
					}
				}
				// если не вышли добавляем месяц
				tmp = time.Date(tmp.Year(), tmp.Month()+1, 1, 0, 0, 0, 0, time.UTC)
			}
			// если первый проход, то присваиваем первое значение и пропускаем цикл
			if first {
				nearestDate = tmp
				first = false
				continue
			}
			// проверяем если tmp раньше , то делаем nearestDate = tmp
			if after(nearestDate, tmp) {
				nearestDate = tmp
			}
		}
		return nearestDate.Format(TIMEFORMAT), nil
	}
	// если месяцы указаны
	if len(rule) == 3 {
		nearestDate := now
		first := true
		// начинаем цикл по месяцам
		for _, m := range strings.Split(rule[2], ",") {
			month, err := strconv.Atoi(m)
			if err != nil {
				return "", fmt.Errorf("cant convert month to int %w", err)
			}
			if month > 12 || month < 1 {
				return "", errors.New("invalid month number")
			}
			// вложенный цикл по дням
			for _, d := range strings.Split(rule[1], ",") {
				day, err := strconv.Atoi(d)
				if err != nil {
					return "", fmt.Errorf("cant convert day to int %w", err)
				}
				if day > 31 || day < -2 {
					return "", errors.New("invalid day number")
				}
				tmp := time.Date(now.Year(), time.Month(month), 1, 0, 0, 0, 0, time.UTC)
				// проверяем можно ли создать такую дату
				if isValidDate(tmp, day) {
					if day < 0 {
						day = daysInMonth(tmp.Year(), tmp.Month()) + day + 1
					}
					tmp = time.Date(tmp.Year(), tmp.Month(), day, 0, 0, 0, 0, time.UTC)
					// если дата меньше now или равна, то добавлем год
					if after(now, tmp) || equal(now, tmp) {
						tmp = time.Date(tmp.Year()+1, tmp.Month(), day, 0, 0, 0, 0, time.UTC)
					}
					// иначе пропускаем цикл, если создать такую дату нельзя
				} else {
					continue
				}
				// если первый раз
				if first {
					nearestDate = tmp
					first = false
					continue
				}
				if after(nearestDate, tmp) && after(tmp, now) {
					nearestDate = tmp
				}
			}
		}
		return nearestDate.Format(TIMEFORMAT), nil
	}
	return "", nil
}

/*
Функция для поиска след. даты по правилу y (каждый год)
*/
func yRule(dstart, now time.Time) string {
	nextDate := dstart
	//в цикле добавляется год
	for {
		nextDate = nextDate.AddDate(1, 0, 0)
		if after(nextDate, now) {
			break
		}
	}
	return nextDate.Format(TIMEFORMAT)
}
