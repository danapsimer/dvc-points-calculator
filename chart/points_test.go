package chart

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestLoadPointChart(t *testing.T) {
	f, err := os.Open("./sample.json")
	if assert.NoError(t, err, "error reading sample.json") {
		pc, err := LoadPointChart(f)
		if assert.NoError(t, err, "error parsing sample.json") {
			assert.Equal(t, "Saratoga Springs", pc.Resort)
			assert.Equal(t, 2022, pc.Year)
			assert.Equal(t, "Deluxe Studio (Standard)", pc.RoomTypes[0].Name)
			assert.Equal(t, 0, pc.RoomTypes[0].Bedrooms)
			assert.Equal(t, 4, pc.RoomTypes[0].Sleeps)
			expectedStartTime, _ := time.Parse("01-02-2006", "09-01-2022")
			expectedEndTime, _ := time.Parse("01-02-2006", "09-30-2022")
			assert.Equal(t, expectedStartTime, pc.Tiers[0].DateRanges[0].StartDate)
			assert.Equal(t, expectedEndTime, pc.Tiers[0].DateRanges[0].EndDate)
			assert.Equal(t, 10, pc.Tiers[0].RoomTypePoints["Deluxe Studio (Standard)"].Weekday)
			assert.Equal(t, 14, pc.Tiers[0].RoomTypePoints["Deluxe Studio (Standard)"].Weekend)
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
		[]string{"Deluxe Studio (Standard)", "One Bedroom Villa (Standard)"},
		map[string]int{
			"Deluxe Studio (Standard)":     14,
			"One Bedroom Villa (Standard)": 29,
		},
	},
}

func TestPointChart_GetPointsForDay(t *testing.T) {
	f, err := os.Open("./sample.json")
	if assert.NoError(t, err, "error reading sample.json") {
		pc, err := LoadPointChart(f)
		if assert.NoError(t, err) {
			for _, s := range GetPointsForDayScenarios {
				t.Run(fmt.Sprintf("%v - %v", s.Date, s.RoomTypes), func(t *testing.T) {
					ptmap, err := pc.GetPointsForDay(s.Date, s.RoomTypes)
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
