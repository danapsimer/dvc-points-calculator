package main

import (
	"context"
	"fmt"
	"github.com/danapsimer/dvc-points-calculator/chart"
	"github.com/danapsimer/dvc-points-calculator/db"
	"github.com/danapsimer/dvc-points-calculator/load/dvcrental"
	"github.com/stretchr/testify/assert"
	"log"
	"strconv"
	"strings"
	"testing"
)

func TestValidateVsDVCRental(t *testing.T) {
	err := db.InitDatastore("dvc-points-calculator-qa", "../google-credentials.json")
	if !assert.NoError(t, err) {
		return
	}
	seasons, resorts, err := dvcrental.LoadSeasonsAndResorts("./dvcrental/seasons.json", "./dvcrental/resorts.json")
	if !assert.NoError(t, err) {
		return
	}
	for _, resortCode := range chart.Resorts {
		resortName := dvcrental.ResortCodeToDVCRentalResortName[resortCode]
		resort := resorts[resortName]
		for seasonId, seasonPrices := range resort.Prices {
			for yearStr, yearPrices := range seasonPrices {
				year, err := strconv.Atoi(yearStr)
				if assert.NoError(t, err) && year >= 2022 {
					ctx := context.Background()
					pc, err := db.GetPointChart(ctx, resortCode, year)
					if err != nil && strings.HasPrefix(err.Error(), "chart not found") {
						log.Printf("no chart found for %s/%d", resortName, year)
					} else if assert.NoError(t, err) {
						season := seasons[seasonId][yearStr]
						if assert.NotNil(t, season) {
							for tierIndex, tier := range pc.Tiers {
								tierRanges := dvcrental.SortAndMergeDateRange(tier.DateRanges)
								var tierName string
								var dvcRentalTierRange dvcrental.SeasonTier
								found := false
								for tierName, dvcRentalTierRange = range season {
									ranges, _ := dvcRentalTierRange.GetRanges()
									if tierRanges.Equal(ranges) {
										found = true
										break
									}
								}
								assert.True(t, found,
									fmt.Sprintf("no matching tier found: seasonId = %s, resortCode=%s, tierIndex = %d", seasonId, resortCode, tierIndex))
								for viewType, viewPrices := range yearPrices {
									tierPrices, ok := viewPrices[tierName]
									if assert.True(t, ok) {
										for roomType, pointStr := range tierPrices.FriSat {
											expectedPoints, _ := strconv.Atoi(pointStr)
											if expectedPoints > 0 {
												roomCode := dvcrental.DVCRentalRoomTypesToCode[resortCode].ViewRoomMap[viewType][roomType]
												points, ok := tier.RoomTypePoints[roomCode]
												if assert.Truef(t, ok, "resortCode = %s, year = %d, tierDates = %s, viewType = %s, roomType = %s, roomCode = %s",
													resortCode, year, tierRanges, viewType, roomType, roomCode) {
													weekendPoints := points.Weekend
													assert.Equalf(t, expectedPoints,
														weekendPoints,
														"resortCode = %s, year = %d, tierDates = %s, viewType = %s, roomType = %s, roomCode = %s",
														resortCode, year, tierRanges, viewType, roomType, roomCode)
												}
											}
										}
										for roomType, pointStr := range tierPrices.SunThur {
											expectedPoints, _ := strconv.Atoi(pointStr)
											if expectedPoints > 0 {
												roomCode := dvcrental.DVCRentalRoomTypesToCode[resortCode].ViewRoomMap[viewType][roomType]
												points, ok := tier.RoomTypePoints[roomCode]
												if assert.Truef(t, ok, "resortCode = %s, year = %d, tierDates = %s, viewType = %s, roomType = %s, roomCode = %s",
													resortCode, year, tierRanges, viewType, roomType, roomCode) {
													assert.Equalf(t, expectedPoints,
														points.Weekday,
														"resortCode = %s, year = %d, tierDates = %s, viewType = %s, roomType = %s, roomCode = %s",
														resortCode, year, tierRanges, viewType, roomType, roomCode)
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}
