package goa

import (
	"context"
	"github.com/danapsimer/dvc-points-calculator/api/goa/gen/points"
	"github.com/danapsimer/dvc-points-calculator/chart"
	"github.com/danapsimer/dvc-points-calculator/model"
	"time"
)

func (s *Points) QueryStay(ctx context.Context, stay *points.Stay) (res *points.StayResult, err error) {
	from, err := time.Parse("2006-01-02", stay.From)
	if err != nil {
		err = points.MakeInvalidInput(err)
		return
	}
	to, err := time.Parse("2006-01-02", stay.To)
	if err != nil {
		err = points.MakeInvalidInput(err)
		return
	}
	stayRequest := &model.Stay{
		From:           from,
		To:             to,
		IncludeResorts: stay.IncludeResorts,
		ExcludeResorts: stay.ExcludeResorts,
		MinSleeps:      stay.MinSleeps,
		MaxSleeps:      stay.MaxSleeps,
		MinBedrooms:    stay.MinBedrooms,
		MaxBedrooms:    stay.MaxBedrooms,
		MinBeds:        stay.MinBeds,
		MaxBeds:        stay.MaxBeds,
	}
	stayResult, err := chart.StayQuery(ctx, stayRequest)
	if err != nil {
		return
	}
	res = &points.StayResult{
		From:           stayResult.From.Format("2006-01-02"),
		To:             stayResult.To.Format("2006-01-02"),
		IncludeResorts: stayResult.IncludeResorts,
		ExcludeResorts: stayResult.ExcludeResorts,
		MinSleeps:      stayResult.MinSleeps,
		MaxSleeps:      stayResult.MaxSleeps,
		MinBedrooms:    stayResult.MinBedrooms,
		MaxBedrooms:    stayResult.MaxBedrooms,
		MinBeds:        stayResult.MinBeds,
		MaxBeds:        stayResult.MaxBeds,
		Rooms:          stayResult.Rooms,
	}
	return
}
