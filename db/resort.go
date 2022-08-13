package db

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/danapsimer/dvc-points-calculator/model"
)

type resortRecord struct {
	ResortCode string `datastore:"resortCode"`
	ResortName string `datastore:"resortName"`
}

func GetResortList(ctx context.Context) ([]*model.Resort, error) {
	query := datastore.NewQuery("PointChart").
		DistinctOn("resortCode", "resortName").
		Project("resortCode", "resortName")
	resortRecords := make([]*resortRecord, 0)
	_, err := dataStoreClient.GetAll(ctx, query, &resortRecords)
	if err != nil {
		return nil, fmt.Errorf("resort list query failed: %s", err.Error())
	}
	result := make([]*model.Resort, 0, len(resortRecords))
	for _, rr := range resortRecords {
		result = append(result, &model.Resort{
			Code: rr.ResortCode,
			Name: rr.ResortName,
		})
	}
	return result, nil
}

func GetResort(ctx context.Context, resortCode string) (*model.Resort, error) {
	query := datastore.NewQuery("PointChart").
		FilterField("resortCode", "=", resortCode).
		DistinctOn("resortCode", "resortName").
		Project("resortCode", "resortName")
	resortRecords := make([]*resortRecord, 0)
	_, err := dataStoreClient.GetAll(ctx, query, &resortRecords)
	if err != nil {
		return nil, fmt.Errorf("resort query failed: %s", err.Error())
	}
	if len(resortRecords) > 1 {
		return nil, fmt.Errorf("multiple resort records found for %s: %+v", resortCode, resortRecords)
	}
	if len(resortRecords) == 0 {
		return nil, nil
	}
	return &model.Resort{
		Code: resortRecords[0].ResortCode,
		Name: resortRecords[0].ResortName,
	}, nil
}

func UpdateResort(ctx context.Context, resort *model.Resort) (*model.Resort, error) {
	_, err := dataStoreClient.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		charts, err := loadPointChartsByCode(ctx, resort.Code)
		if err != nil {
			return fmt.Errorf("error retrieving point charts for %s: %s", resort.Code, err.Error())
		}
		for _, chart := range charts {
			chart.Resort = resort.Name
			err := SavePointChart(ctx, chart)
			if err != nil {
				rerr := tx.Rollback()
				if rerr != nil {
					return fmt.Errorf("error rolling back resort update transaction for %s: original error = %s, rollback error = %s",
						resort.Code, err.Error(), rerr.Error())
				}
				return err
			}
		}
		return nil
	})
	return resort, err
}
