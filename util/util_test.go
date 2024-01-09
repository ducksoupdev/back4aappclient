package util

import "testing"

func TestToBack4AppDate(t *testing.T) {
	date := ToBack4AppDate("date")
	if date.Type != "Date" {
		t.Errorf("expected %s, got %s", "Date", date.Type)
	}
	if date.Iso != "date" {
		t.Errorf("expected %s, got %s", "date", date.Iso)
	}
}

func TestParseBack4AppDate(t *testing.T) {
	date := ParseBack4AppDate(map[string]interface{}{"iso": "date"})
	if date.Type != "Date" {
		t.Errorf("expected %s, got %s", "Date", date.Type)
	}
	if date.Iso != "date" {
		t.Errorf("expected %s, got %s", "date", date.Iso)
	}
}

func TestBack4AppDateToIsoString(t *testing.T) {
	date := Back4AppDateToIsoString(Back4AppDate{Type: "Date", Iso: "date"})
	if date != "date" {
		t.Errorf("expected %s, got %s", "date", date)
	}
}

func TestBack4AppDateToTime(t *testing.T) {
	date, _ := Back4AppDateToTime(Back4AppDate{Type: "Date", Iso: "2020-01-01T00:00:00.000Z"})
	if date.Year() != 2020 {
		t.Errorf("expected %d, got %d", 2020, date.Year())
	}
	if date.Month() != 1 {
		t.Errorf("expected %d, got %d", 1, date.Month())
	}
	if date.Day() != 1 {
		t.Errorf("expected %d, got %d", 1, date.Day())
	}
	if date.Hour() != 0 {
		t.Errorf("expected %d, got %d", 0, date.Hour())
	}
	if date.Minute() != 0 {
		t.Errorf("expected %d, got %d", 0, date.Minute())
	}
	if date.Second() != 0 {
		t.Errorf("expected %d, got %d", 0, date.Second())
	}
	if date.Nanosecond() != 0 {
		t.Errorf("expected %d, got %d", 0, date.Nanosecond())
	}
}
