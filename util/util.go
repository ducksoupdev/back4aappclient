package util

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
