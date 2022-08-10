package model

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestLoadPointChart(t *testing.T) {
	f, err := os.Open("../load/ssr/2022.json")
	if assert.NoError(t, err, "error reading /load/ssr/2022.json") {
		pc, err := ReadPointChart(f)
		if assert.NoError(t, err, "error parsing /load/ssr/2022.json") {
			assert.Equal(t, "Disney's Saratoga Springs Resort & Spa", pc.Resort)
			assert.Equal(t, 2022, pc.Year)
			assert.Equal(t, "Deluxe Studio (Standard)", pc.RoomTypes[0].Name)
			assert.Equal(t, 0, pc.RoomTypes[0].Bedrooms)
			assert.Equal(t, 4, pc.RoomTypes[0].Sleeps)
			expectedStartTime, err := ParseTierDate("09-01")
			assert.NoError(t, err)
			expectedEndTime, err := ParseTierDate("09-30")
			assert.NoError(t, err)
			assert.Equal(t, expectedStartTime, pc.Tiers[0].DateRanges[0].StartDate)
			assert.Equal(t, expectedEndTime, pc.Tiers[0].DateRanges[0].EndDate)
			assert.Equal(t, 10, pc.Tiers[0].RoomTypePoints["dss"].Weekday)
			assert.Equal(t, 14, pc.Tiers[0].RoomTypePoints["dss"].Weekend)
		}
	}
}

var GetPointsForDayScenarios = []struct {
	Date           time.Time
	RoomTypes      []string
	ExpectedPoints map[string]int
}{
	{
		time.Date(2022, time.February, 15, 0, 0, 0, 0, time.UTC),
		[]string{"dss", "1bs"},
		map[string]int{
			"dss": 14,
			"1bs": 29,
		},
	},
}

func TestPointChart_GetPointsForDay(t *testing.T) {
	f, err := os.Open("../load/ssr/2022.json")
	if assert.NoError(t, err, "error reading /load/ssr/2022.json") {
		pc, err := ReadPointChart(f)
		if assert.NoError(t, err) {
			for _, s := range GetPointsForDayScenarios {
				t.Run(fmt.Sprintf("%v - %v", s.Date, s.RoomTypes), func(t *testing.T) {
					ptmap, err := pc.GetPointsForDay(s.Date, s.RoomTypes...)
					if assert.NoError(t, err) {
						assert.Equal(t, len(s.RoomTypes), len(ptmap))
						for _, rt := range s.RoomTypes {
							assert.Equal(t, s.ExpectedPoints[rt], ptmap[rt])
						}
					}
				})
			}
		}
	}
}

func TestTierDate_MarshalJSON(t *testing.T) {
	type Test struct {
		Date TierDate `json:"date"`
	}
	scenarios := []struct {
		input    Test
		expected any
	}{
		{
			Test{TierDate(time.Date(2022, 01, 01, 0, 0, 0, 0, time.UTC))},
			"{\"date\":\"01-01\"}",
		},
		{
			Test{TierDate(time.Date(2023, 12, 15, 0, 0, 0, 0, time.UTC))},
			"{\"date\":\"12-15\"}",
		},
	}
	for idx, scenario := range scenarios {
		t.Run(fmt.Sprintf("Scenario %d", idx), func(t *testing.T) {
			marshalled, err := json.Marshal(scenario.input)
			if assert.NoError(t, err) {
				assert.Equal(t, scenario.expected, string(marshalled))
			}
		})
	}
}

func TestTierDate_UnmarshalJSON(t *testing.T) {
	type Test struct {
		Date TierDate `json:"date"`
	}
	scenarios := []struct {
		expected any
		input    string
	}{
		{
			Test{TierDate(time.Date(2022, 01, 01, 0, 0, 0, 0, time.UTC))},
			"{\"date\":\"01-01\"}",
		},
		{
			Test{TierDate(time.Date(2022, 12, 15, 0, 0, 0, 0, time.UTC))},
			"{\"date\":\"12-15\"}",
		},
		{
			fmt.Errorf("malformed json string: 1215"),
			"{\"date\":1215}",
		},
	}
	for idx, scenario := range scenarios {
		t.Run(fmt.Sprintf("Scenario %d", idx), func(t *testing.T) {
			var unmarshalled Test
			err := json.Unmarshal([]byte(scenario.input), &unmarshalled)
			if expectedErr, ok := scenario.expected.(error); ok {
				assert.Equal(t, expectedErr, err)
			} else if assert.NoError(t, err) {
				assert.Equal(t, scenario.expected, unmarshalled)
			}
		})
	}
}

func TestTierDate_Day(t *testing.T) {
	now := time.Now()
	assert.Equal(t, now.Day(), TierDate(now).Day())
}

func TestTierDate_Month(t *testing.T) {
	now := time.Now()
	assert.Equal(t, now.Month(), TierDate(now).Month())
}

func TestParseTierDate(t *testing.T) {
	scenarios := []struct {
		expected any
		input    string
	}{
		{
			TierDate(time.Date(2022, 01, 01, 0, 0, 0, 0, time.UTC)),
			"01-01",
		},
		{
			TierDate(time.Date(2022, 12, 15, 0, 0, 0, 0, time.UTC)),
			"12-15",
		},
		{
			fmt.Errorf("malformed tier date: AA-15 - expected a numeric month component: AA"),
			"AA-15",
		},
		{
			fmt.Errorf("malformed tier date: 12-AA - expected a numeric day component: AA"),
			"12-AA",
		},
		{
			fmt.Errorf("malformed tier date: 1215 - expected a month and day seperated by '-'"),
			"1215",
		},
	}
	for idx, scenario := range scenarios {
		t.Run(fmt.Sprintf("Scenario %d", idx), func(t *testing.T) {
			unmarshalled, err := ParseTierDate(scenario.input)
			if expectedErr, ok := scenario.expected.(error); ok {
				assert.Equal(t, expectedErr, err)
			} else if assert.NoError(t, err) {
				assert.Equal(t, scenario.expected, unmarshalled)
			}
		})
	}
}
