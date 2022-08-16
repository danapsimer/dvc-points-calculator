package goa

import (
	"context"
	"fmt"
	dvcpointscalculator "github.com/danapsimer/dvc-points-calculator/api/goa/gen/dvc_points_calculator"
	"github.com/danapsimer/dvc-points-calculator/db"
)

func (s *DVCPointsCalculatorService) GetResorts(ctx context.Context) (res dvcpointscalculator.DvcpointcalculatorResortCollection, view string, err error) {
	view = "resortOnly"
	resorts, err := db.GetResortList(ctx)
	if err != nil {
		return
	}
	res = make([]*dvcpointscalculator.DvcpointcalculatorResort, 0, len(resorts))
	for _, resort := range resorts {
		newResort := new(dvcpointscalculator.DvcpointcalculatorResort)
		newResort.CreateFromResort(resort)
		res = append(res, newResort)
	}
	return
}

func (s *DVCPointsCalculatorService) GetResort(ctx context.Context, payload *dvcpointscalculator.GetResortPayload) (res *dvcpointscalculator.DvcpointcalculatorResort, view string, err error) {
	view = "resortOnly"
	res = new(dvcpointscalculator.DvcpointcalculatorResort)
	resort, err := db.GetResort(ctx, string(payload.ResortCode))
	if err != nil {
		return
	} else if resort != nil {
		res.CreateFromResort(resort)
	} else {
		err = dvcpointscalculator.MakeNotFound(fmt.Errorf("no such resort: %s", payload.ResortCode))
	}
	return
}
