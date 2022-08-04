package chart

import (
	"context"
	"fmt"
	"time"
)

var (
	ErrorChartNotFound = "chart not found"
)

type Stay struct {
	From time.Time `json:"from" uri:"from" binding:"required,ltefield=To" time_format:"2006-01-02" time_utc:"1"`
	To   time.Time `json:"to" uri:"to" binding:"required" time_format:"2006-01-02" time_utc:"1"`
}

type StayResult struct {
	Stay
	Rooms map[string]map[string]int `json:"rooms"`
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

func StayQuery(ctx context.Context, stay *Stay) (*StayResult, error) {
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
		result1, err := StayQuery(ctx, query1)
		if err != nil {
			return nil, err
		}
		result2, err := StayQuery(ctx, query2)
		if err != nil {
			return nil, err
		}
		result.mergeResults(result1, result2)
	} else {
		for _, resortCode := range Resorts {
			chart, err := LoadPointChartByCodeAndYear(ctx, resortCode, stay.From.Year())
			if err != nil {
				return nil, err
			} else {
				result.Rooms[resortCode], err = chart.GetPointsForStay(stay)
				if err != nil {
					return nil, fmt.Errorf("error calculating points for a stay: %+v - %+v", stay, err)
				}
			}
		}
	}
	return result, nil
}
