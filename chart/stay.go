package chart

import (
	"context"
	"fmt"
	"github.com/danapsimer/dvc-points-calculator/db"
	"github.com/danapsimer/dvc-points-calculator/model"
	"github.com/danapsimer/dvc-points-calculator/util"
	"time"
)

var (
	ErrorChartNotFound = "chart not found"
	Resorts            = []string{
		"akv",
		"aul",
		"blt",
		"bcv",
		"bwv",
		"brv",
		"ccv",
		"gcv",
		"gfv",
		"hhr",
		"okw",
		"pvv",
		"riv",
		"ssr",
		"vbr",
	}
)

func NewStayResult(stay *model.Stay) *model.StayResult {
	return &model.StayResult{Stay: *stay, Rooms: make(map[string]map[string]int)}
}

func StayQuery(ctx context.Context, stay *model.Stay) (*model.StayResult, error) {
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
		query1 := *stay
		query1.To = time.Date(stay.To.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		query2 := *stay
		query2.From = time.Date(stay.To.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		result1, err := StayQuery(ctx, &query1)
		if err != nil {
			return nil, err
		}
		result2, err := StayQuery(ctx, &query2)
		if err != nil {
			return nil, err
		}
		result.MergeResults(result1, result2)
	} else {
		resortsToSearch := make([]string, len(Resorts))
		copy(resortsToSearch, Resorts)
		for idx := len(resortsToSearch) - 1; idx >= 0; idx-- {
			if (stay.IncludeResorts != nil && !util.Contains(stay.IncludeResorts, resortsToSearch[idx])) ||
				(stay.ExcludeResorts != nil && util.Contains(stay.ExcludeResorts, resortsToSearch[idx])) {
				if idx < len(resortsToSearch)-1 {
					resortsToSearch = append(resortsToSearch[:idx], resortsToSearch[idx+1:]...)
				} else {
					resortsToSearch = resortsToSearch[:idx]
				}
			}
		}
		for _, resortCode := range resortsToSearch {
			chart, err := db.GetPointChart(ctx, resortCode, stay.From.Year())
			if err != nil {
				return nil, err
			} else if chart != nil {
				result.Rooms[resortCode], err = chart.GetPointsForStay(stay)
				if err != nil {
					return nil, fmt.Errorf("error calculating points for a stay: %+v - %+v", stay, err)
				}
			}
		}
	}
	return result, nil
}
