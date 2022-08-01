package chart

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type DateRange struct {
	StartDate time.Time
	EndDate   time.Time
}

// Contains date must be truncated to midnight in UTC.
func (dr *DateRange) Contains(date time.Time) bool {
	return dr.StartDate.Equal(date) || (dr.StartDate.Before(date) && (dr.EndDate.Equal(date) || dr.EndDate.After(date)))
}

type RoomType struct {
	Name     string
	Sleeps   int
	Bedrooms int
}

type Points struct {
	Weekday, Weekend int
}

type Tier struct {
	DateRanges     []DateRange
	RoomTypePoints map[string]Points
}

type PointChart struct {
	Resort    string
	Year      int
	RoomTypes []RoomType
	Tiers     []Tier
}

var PointCharts map[string]map[int]*PointChart

func LoadPointChart(in io.Reader) (*PointChart, error) {
	decoder := json.NewDecoder(in)
	pc := &PointChart{}
	err := decoder.Decode(pc)
	if err != nil {
		return nil, err
	}
	return pc, nil
}

func (pc *PointChart) GetPointsForDay(date time.Time, roomTypes ...string) (map[string]int, error) {
	if pc.Year != date.Year() {
		return nil, fmt.Errorf("date out of range for point chart: got %d, but expected %d", date.Year(), pc.Year)
	}
	weekday := date.Weekday()
	isWeekend := weekday == time.Friday || weekday == time.Saturday
	if roomTypes == nil || len(roomTypes) == 0 {
		roomTypes = make([]string, 0, len(pc.RoomTypes))
		for _, rt := range pc.RoomTypes {
			roomTypes = append(roomTypes, rt.Name)
		}
	}
	for _, t := range pc.Tiers {
		for _, dr := range t.DateRanges {
			if dr.Contains(date) {
				pointsMap := make(map[string]int)
				for _, rt := range roomTypes {
					if isWeekend {
						pointsMap[rt] = t.RoomTypePoints[rt].Weekend
					} else {
						pointsMap[rt] = t.RoomTypePoints[rt].Weekday
					}
				}
				return pointsMap, nil
			}
		}
	}
	return nil, fmt.Errorf("date didn't match any tier: %v", date)
}

// GetPointsForStay calculates the points for all room types in the chart for the given stay.
// from should be truncated to midnight in UTC.
// to should be truncated to midnight in UTC
// Returns a map containing room type name as the key and points as the value.
func (pc *PointChart) GetPointsForStay(stay *Stay) (map[string]int, error) {
	points := make(map[string]int)
	for date := stay.From; date.Before(stay.To); date = date.Add(time.Hour * 24) {
		dayPoints, err := pc.GetPointsForDay(date)
		if err != nil {
			return nil, err
		}
		for rt, pts := range dayPoints {
			points[rt] += pts
		}
	}
	return points, nil
}
