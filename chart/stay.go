package chart

import (
	"fmt"
	"time"
)

type Stay struct {
	From time.Time `json:"from" uri:"from" binding:"required"`
	To   time.Time `json:"to" uri:"to" binding:"required"`
}

type StayResult struct {
	Stay
	Rooms map[string]map[string]int
}

func NewStayResult(stay *Stay) *StayResult {
	return &StayResult{*stay, make(map[string]map[string]int)}
}

func (sr *StayResult) addPoints(resort, roomType string, points int) {
	if sr.Rooms[resort] == nil {
		sr.Rooms[resort] = make(map[string]int)
	}
	sr.Rooms[resort][roomType] += points
}

func (sr *StayResult) mergeResults(results ...*StayResult) {
	for _, result := range results {
		for resort, resortPoints := range result.Rooms {
			for roomType, roomTypePoints := range resortPoints {
				sr.addPoints(resort, roomType, roomTypePoints)
			}
		}
	}
}

func StayQuery(stay *Stay) (*StayResult, error) {
	result := NewStayResult(stay)
	if !stay.To.After(stay.From) {
		return nil, fmt.Errorf("to date must be strictly after from date: %v - %v is not a valid date range",
			stay.From.Format("2006-01-02"), stay.To.Format("2006-01-02"))
	}
	if stay.To.Sub(stay.From) > time.Hour*24*365 {
		return nil, fmt.Errorf("to cannot be more than 365 days after from: %s - %s is not a valid date range",
			stay.From.Format("2006-01-02"), stay.To.Format("2006-01-02"))
	}

	if stay.From.Year() != stay.To.Year() && stay.To.YearDay() > 1 {
		query1 := &Stay{stay.From, time.Date(stay.To.Year(), 1, 1, 0, 0, 0, 0, time.UTC)}
		query2 := &Stay{time.Date(stay.To.Year(), 1, 1, 0, 0, 0, 0, time.UTC), stay.To}
		result1, err := StayQuery(query1)
		if err != nil {
			return nil, err
		}
		result2, err := StayQuery(query2)
		if err != nil {
			return nil, err
		}
		result.mergeResults(result1, result2)
	} else {
		for resort, charts := range PointCharts {
			chart := charts[stay.From.Year()]
			var err error
			result.Rooms[resort], err = chart.GetPointsForStay(stay)
			if err != nil {
				return nil, err
			}
		}
	}
	return result, nil
}
