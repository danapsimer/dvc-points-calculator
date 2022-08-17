package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type TierDate string

func ParseTierDate(dateStr string) (TierDate, error) {
	var tierDate TierDate

	parts := strings.Split(dateStr, "-")
	if len(parts) != 2 {
		return tierDate, fmt.Errorf("malformed tier date: %s - expected a month and day seperated by '-'", dateStr)
	}
	month, err := strconv.Atoi(parts[0])
	if err != nil {
		return tierDate, fmt.Errorf("malformed tier date: %s - expected a numeric month component: %s", dateStr, parts[0])
	}
	day, err := strconv.Atoi(parts[1])
	if err != nil {
		return tierDate, fmt.Errorf("malformed tier date: %s - expected a numeric day component: %s", dateStr, parts[1])
	}
	tierDate = TierDate(fmt.Sprintf("%d-%d", month, day))
	return tierDate, nil
}

func (t TierDate) MarshalJSON() ([]byte, error) {
	return []byte("\"" + t + "\""), nil
}

func (t *TierDate) UnmarshalJSON(date []byte) error {
	dateStr := string(date)
	if !strings.HasPrefix(dateStr, "\"") || !strings.HasSuffix(dateStr, "\"") {
		return fmt.Errorf("malformed json string: %s", dateStr)
	}
	dateStr = dateStr[1 : len(dateStr)-1]
	d, err := ParseTierDate(dateStr)
	if err != nil {
		return err
	}
	*t = d
	return nil
}

func (t TierDate) Month() time.Month {
	parts := strings.Split(string(t), "-")
	month, _ := strconv.Atoi(parts[0])
	return time.Month(month)
}

func (t TierDate) Day() int {
	parts := strings.Split(string(t), "-")
	day, _ := strconv.Atoi(parts[1])
	return day
}

func (t TierDate) Before(s TierDate) bool {
	tmonth := t.Month()
	smonth := s.Month()
	return tmonth < smonth || (tmonth == smonth && t.Day() < s.Day())
}

func (t TierDate) AsTime() time.Time {
	return time.Date(time.Now().Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func FromTime(t time.Time) TierDate {
	return TierDate(fmt.Sprintf("%d-%d", t.Month(), t.Day()))
}

func (t TierDate) Add(d time.Duration) TierDate {
	return FromTime(t.AsTime().Add(d).Truncate(time.Hour * 24))
}

func (t TierDate) String() string {
	return string(t)
}

type DateRange struct {
	StartDate TierDate `json:"startDate"`
	EndDate   TierDate `json:"endDate"`
}

func (dr *DateRange) String() string {
	return fmt.Sprintf("%s - %s", dr.StartDate, dr.EndDate)
}

func (dr *DateRange) Equal(ddr *DateRange) bool {
	return dr.StartDate == ddr.StartDate && dr.EndDate == ddr.EndDate
}

// Contains date must be truncated to midnight in UTC.
func (dr *DateRange) Contains(date time.Time) bool {
	from := time.Date(date.Year(), dr.StartDate.Month(), dr.StartDate.Day(), 0, 0, 0, 0, time.UTC)
	to := time.Date(date.Year(), dr.EndDate.Month(), dr.EndDate.Day(), 0, 0, 0, 0, time.UTC)
	return from.Equal(date) || (from.Before(date) && (to.Equal(date) || to.After(date)))
}
