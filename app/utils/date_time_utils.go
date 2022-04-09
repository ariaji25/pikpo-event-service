package utils

import (
	"time"
)

var dateLayoutFormat = "2006-01-02"
var timeLayoutFormat = "15:04:05"

func ToDateTime(s string) (*time.Time, error) {
	dateTime, err := time.Parse(dateLayoutFormat, s)
	if err != nil {
		return nil, err
	}
	return &dateTime, nil
}

func ToTime(s string) (*time.Time, error) {
	time, err := time.Parse(timeLayoutFormat, s)
	if err != nil {
		return nil, err
	}
	return &time, nil
}
