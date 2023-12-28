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
