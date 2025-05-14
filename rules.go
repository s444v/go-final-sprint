package main

import (
	"fmt"
	"time"
)

const plusYear = 31_536_000

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	timestamp := now.Unix()
	timestamp += plusYear
	t := time.Unix(timestamp, 0)
	formatted := t.Format("20060102")
	return formatted, nil
}

func main() {
	now := time.Date(2024, 01, 26, 0, 0, 0, 0, time.UTC)
	res, err := NextDate(now, "20240229", "y")
	fmt.Println(res, err)
}
