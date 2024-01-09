package util

import "time"

type Back4AppDate struct {
	Type string `json:"__type"`
	Iso  string `json:"iso"`
}

func ToBack4AppDate(date string) Back4AppDate {
	return Back4AppDate{
		Type: "Date",
		Iso:  date,
	}
}

func ParseBack4AppDate(date map[string]interface{}) Back4AppDate {
	return ToBack4AppDate(date["iso"].(string))
}

func Back4AppDateToIsoString(date Back4AppDate) string {
	return date.Iso
}

func Back4AppDateToTime(date Back4AppDate) (time.Time, error) {
	return time.Parse(time.RFC3339, date.Iso)
}
