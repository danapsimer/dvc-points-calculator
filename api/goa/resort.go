package goa

import (
	"context"
	"fmt"
	"github.com/danapsimer/dvc-points-calculator/api/goa/gen/points"
	"github.com/danapsimer/dvc-points-calculator/chart"
	"github.com/danapsimer/dvc-points-calculator/db"
	"strings"
)

func (s *Points) GetResorts(ctx context.Context) (res points.ResortResultCollection, view string, err error) {
	view = "resortOnly"
	resorts, err := db.GetResortList(ctx)
	if err != nil {
		return
	}
	res = make([]*points.ResortResult, 0, len(resorts))
	for _, resort := range resorts {
		newResort := new(points.ResortResult)
		newResort.CreateFromResort(resort)
		res = append(res, newResort)
	}
	return
}

func (s *Points) GetResort(ctx context.Context, payload *points.GetResortPayload) (res *points.ResortResult, err error) {
	resort, err := db.GetResort(ctx, string(payload.ResortCode))
	if err != nil {
		return
	} else if resort != nil {
		res = new(points.ResortResult)
		res.CreateFromResort(resort)
	} else {
		err = points.MakeNotFound(fmt.Errorf("no such resort: %s", payload.ResortCode))
	}
	return
}

func (s *Points) PutResort(ctx context.Context, payload *points.PutResortPayload) (res *points.ResortResult, err error) {
	resort, err := db.GetResort(ctx, string(payload.ResortCode))
	if err != nil {
		return
	} else if resort != nil {
		resort.Name = payload.Name
		if resort, err = db.UpdateResort(ctx, resort); err != nil {
			return
		}
		res = new(points.ResortResult)
		res.CreateFromResort(resort)
	} else {
		err = points.MakeNotFound(fmt.Errorf("no such resort: %s", payload.ResortCode))
	}
	return
}

func (s *Points) GetResortYear(ctx context.Context, payload *points.GetResortYearPayload) (res *points.ResortYearResult, err error) {
	res = new(points.ResortYearResult)
	pointChart, err := db.GetPointChart(ctx, payload.ResortCode, payload.Year)
	if err != nil {
		if strings.HasPrefix(err.Error(), chart.ErrorChartNotFound) {
			err = points.MakeNotFound(fmt.Errorf("no such resort: %s/%d", payload.ResortCode, payload.Year))
		}
		return
	} else if pointChart != nil {
		res.Code = &payload.ResortCode
		res.Name = &pointChart.Resort
		res.Year = &payload.Year
		res.RoomTypes = make([]*points.RoomType, 0, len(pointChart.RoomTypes))
		for _, rt := range pointChart.RoomTypes {
			outrt := new(points.RoomType)
			outrt.CreateFromRoomType(rt)
			res.RoomTypes = append(res.RoomTypes, outrt)
		}
	}
	return
}
