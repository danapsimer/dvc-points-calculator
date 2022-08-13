package db

import (
	"bytes"
	"cloud.google.com/go/datastore"
	"context"
	"encoding/json"
	"fmt"
	"github.com/danapsimer/dvc-points-calculator/model"
	"google.golang.org/api/option"
)

var (
	dataStoreClient *datastore.Client
	pointCharts     map[string]map[int]*model.PointChart
)

func AddPointChart(c *model.PointChart) {
	if pointCharts == nil {
		pointCharts = make(map[string]map[int]*model.PointChart)
	}
	var resortCharts map[int]*model.PointChart
	resortCharts, ok := pointCharts[c.Code]
	if !ok {
		resortCharts = make(map[int]*model.PointChart)
		pointCharts[c.Code] = resortCharts
	}
	resortCharts[c.Year] = c
}

func GetPointChart(ctx context.Context, code string, year int) (*model.PointChart, error) {
	if pointCharts == nil {
		pointCharts = make(map[string]map[int]*model.PointChart)
	}
	resortCharts, ok := pointCharts[code]
	var chart *model.PointChart
	if ok {
		chart, ok = resortCharts[year]
	} else {
		resortCharts = make(map[int]*model.PointChart)
		pointCharts[code] = resortCharts
	}
	if !ok {
		var err error
		chart, err = loadPointChartByCodeAndYear(ctx, code, year)
		if err != nil {
			return nil, err
		}
		resortCharts[year] = chart
	}
	return chart, nil
}

func InitDatastore(projectId, credentialFile string) (err error) {
	ctx := context.Background()
	dataStoreClient, err = datastore.NewClient(ctx, projectId, option.WithCredentialsFile(credentialFile))
	if err != nil {
		err = fmt.Errorf("error creating new client: %+v", err)
	}
	return
}

type pointChartRecord struct {
	ResortCode string `json:"resortCode" datastore:"resortCode"`
	ResortName string `json:"resortName" datastore:"resortName"`
	Year       int    `json:"year" datastore:"year"`
	Chart      []byte `json:"chart" datastore:"chart,noindex"`
}

func loadPointChartByCodeAndYear(ctx context.Context, resortCode string, year int) (*model.PointChart, error) {
	query := datastore.NewQuery("PointChart").
		FilterField("resortCode", "=", resortCode).
		FilterField("year", "=", year)
	charts := make([]*pointChartRecord, 0)
	_, err := dataStoreClient.GetAll(ctx, query, &charts)
	if err != nil {
		return nil, fmt.Errorf("query failed for %s - %d: %+v", resortCode, year, err)
	}
	if l := len(charts); l > 1 {
		return nil, fmt.Errorf("multiple charts found %s - %d (%d)", resortCode, year, l)
	}
	for _, pointChart := range charts {
		pc, err := model.ReadPointChart(bytes.NewBuffer(pointChart.Chart))
		if err != nil {
			return nil, fmt.Errorf("point chart parsing failed for %s - %d: %+v", resortCode, year, err)
		}
		AddPointChart(pc)
		return pc, nil
	}
	return nil, fmt.Errorf("chart not found %s - %d", resortCode, year)
}

func SavePointChart(ctx context.Context, chart *model.PointChart) error {
	buffer := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buffer)
	encoder.Encode(chart)

	pcr := &pointChartRecord{
		chart.Code,
		chart.Resort,
		chart.Year,
		buffer.Bytes(),
	}
	key := datastore.NameKey("PointChart", fmt.Sprintf("%s_%d", chart.Code, chart.Year), nil)
	_, err := dataStoreClient.Put(ctx, key, pcr)
	return err
}
