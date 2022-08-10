package main

import (
	"context"
	"dvccalc/chart"
	"dvccalc/db"
	"dvccalc/model"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

var filepathRegexp = regexp.MustCompile(".*/([a-z]+)/(\\d{4}).json$")

func contains[E comparable](v []E, find E) bool {
	for _, e := range v {
		if e == find {
			return true
		}
	}
	return false
}

func LoadPointCharts() error {
	ctx := context.Background()
	return filepath.Walk(filepath.Clean("."), func(path string, info fs.FileInfo, err error) error {
		match := filepathRegexp.FindStringSubmatch(path)
		if match != nil && len(match) == 3 && contains(chart.Resorts, match[1]) {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			pc, err := model.ReadPointChart(f)
			if err != nil {
				return fmt.Errorf("error loading %s - %s", path, err.Error())
			}
			err = db.SavePointChart(ctx, pc)
			if err != nil {
				return fmt.Errorf("error saving point chart: %s - %s", path, err.Error())
			}
		}
		return nil
	})
}

func main() {
	err := db.InitDatastore("dvc-points-calculator-qa", "./google-credentials.json")
	if err != nil {
		panic(err)
	}
	err = LoadPointCharts()
	if err != nil {
		panic(err)
	}
}
