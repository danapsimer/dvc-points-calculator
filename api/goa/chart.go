package goa

import (
	"context"
	"github.com/danapsimer/dvc-points-calculator/api/goa/gen/points"
	"github.com/danapsimer/dvc-points-calculator/db"
)

func (s *Points) GetPointChart(ctx context.Context, payload *points.GetPointChartPayload) (res *points.PointChart, err error) {
	res = new(points.PointChart)
	if pointChart, err := db.GetPointChart(ctx, payload.ResortCode, payload.Year); err == nil {
		res.CreateFromPointChart(pointChart)
	}
	return
}
