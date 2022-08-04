package chart

import (
	"cloud.google.com/go/datastore"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/api/option"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
	Name     string `json:"name"`
	Code     string `json:"code"`
	Sleeps   int    `json:"sleeps"`
	Bedrooms int    `json:"bedrooms"`
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

var (
	Resorts = []string{
		"ssr",
		"aul",
	}
	PointCharts     map[string]map[int]*PointChart
	dataStoreClient *datastore.Client
)

func LoadPointChart(in io.Reader) (*PointChart, error) {
	decoder := json.NewDecoder(in)
	pc := &PointChart{}
	err := decoder.Decode(pc)
	if err != nil {
		return nil, fmt.Errorf("error decoding point chart: %s", err)
	}
	return pc, nil
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
	ResortCode string `json:"resortCode"`
	ResortName string `json:"resortName"`
	Year       int    `json:"year"`
	Chart      string `json:"chart"`
}

func LoadPointChartByCodeAndYear(ctx context.Context, resortCode string, year int) (*PointChart, error) {
	if PointCharts == nil {
		PointCharts = make(map[string]map[int]*PointChart)
	}
	if resort, ok := PointCharts[resortCode]; ok {
		if chart, ok := resort[year]; ok {
			return chart, nil
		}
	} else {
		PointCharts[resortCode] = make(map[int]*PointChart)
	}
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
		pc, err := LoadPointChart(strings.NewReader(pointChart.Chart))
		if err != nil {
			return nil, fmt.Errorf("point chart parsing failed for %s - %d: %+v", resortCode, year, err)
		}
		PointCharts[pointChart.ResortCode][year] = pc
		return pc, nil
	}
	return nil, fmt.Errorf("chart not found %s - %d", resortCode, year)
}

func contains[E comparable](v []E, find E) bool {
	for _, e := range v {
		if e == find {
			return true
		}
	}
	return false
}

var filepathRegexp = regexp.MustCompile(".*/([a-z]+)/(\\d{4}).json$")

func LoadPointCharts() error {
	PointCharts = make(map[string]map[int]*PointChart)
	return filepath.Walk(filepath.Clean("."), func(path string, info fs.FileInfo, err error) error {
		match := filepathRegexp.FindStringSubmatch(path)
		if match != nil && len(match) == 3 && contains(Resorts, match[1]) {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			chart, err := LoadPointChart(f)
			if err != nil {
				return fmt.Errorf("error loading %s - %s", path, err)
			}
			resort := PointCharts[chart.Resort]
			if resort == nil {
				resort = make(map[int]*PointChart)
				PointCharts[chart.Resort] = resort
			}
			resort[chart.Year] = chart
		}
		return nil
	})
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
			roomTypes = append(roomTypes, rt.Code)
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
			return nil, fmt.Errorf("error getting points for a date: %+v - %s : %+v", stay, date.String(), err)
		}
		for rt, pts := range dayPoints {
			points[rt] += pts
		}
	}
	return points, nil
}
