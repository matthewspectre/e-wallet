package wallet

import (
	"testing"
	"time"
)

func TestParseDateQuery(t *testing.T) {
	t.Run("empty date is omitted", func(t *testing.T) {
		date, err := parseDateQuery("")
		if err != nil {
			t.Fatalf("parseDateQuery returned an error: %v", err)
		}
		if date != nil {
			t.Fatalf("expected nil date, got %v", date)
		}
	})

	t.Run("valid date", func(t *testing.T) {
		date, err := parseDateQuery("2026-11-01")
		if err != nil {
			t.Fatalf("parseDateQuery returned an error: %v", err)
		}
		if want := time.Date(2026, time.November, 1, 0, 0, 0, 0, time.UTC); date == nil || !date.Equal(want) {
			t.Fatalf("expected date %v, got %v", want, date)
		}
	})

	t.Run("invalid date", func(t *testing.T) {
		if _, err := parseDateQuery("01-11-2026"); err == nil {
			t.Fatal("expected invalid date format to be rejected")
		}
	})
}
