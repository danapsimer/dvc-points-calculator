package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/danapsimer/dvc-points-calculator/chart"
	"github.com/danapsimer/dvc-points-calculator/db"
	"github.com/danapsimer/dvc-points-calculator/load/dvcrental"
	"github.com/danapsimer/dvc-points-calculator/model"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

var (
	GoogleProjectId       string
	GoogleCredentialsFile string
	SeasonsFileName       string
	ResortsFileName       string
	ResortCodes           string
	Command               string
)

func init() {
	flag.StringVar(&Command, "command", "", "-command=<command>  *Required*")
	flag.StringVar(&GoogleProjectId, "projectId", "dvc-points-calculator-qa", "-projectId=<google project name>")
	flag.StringVar(&GoogleCredentialsFile, "credentials", "./google-credentials.json", "-credentials=<google service account credentials>")
	flag.StringVar(&SeasonsFileName, "seasonsFile", "./load/dvcrental/seasons.json", "-seasonsFile=<seasons file name>")
	flag.StringVar(&ResortsFileName, "resortsFile", "./load/dvcrental/resorts.json", "-resortsFile=<resorts file name>")
	flag.StringVar(&ResortCodes, "resortCodes", strings.Join(chart.Resorts, ","), "-resortCodes=<comma seperated list of resort codes>")
}

func main() {
	flag.Parse()

	chart.Resorts = strings.Split(ResortCodes, ",")

	err := db.InitDatastore(GoogleProjectId, GoogleCredentialsFile)
	if err != nil {
		panic(err)
	}
	switch Command {
	case "loadFromLocalFiles":
		err = LoadPointCharts()
	case "loadFromDVCRentalFiles":
		err = dvcrental.ConvertDVCRentalCharts(chart.Resorts, SeasonsFileName, ResortsFileName)
	default:
		flag.Usage()
	}
	if err != nil {
		panic(err)
	}
}
