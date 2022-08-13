package dvcrental

import (
	"bytes"
	"context"
	"dvccalc/db"
	"dvccalc/model"
	"dvccalc/util"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
)

var (
	DVCRentalResortNameToCode = map[string]string{
		"animal_kingdom_villas": "akv",
		"aulani":                "aus",
		"bay_lake_tower":        "blt",
		"beach_club_villas":     "bcv",
		"boardwalk_villas":      "bwv",
		"boulder_ridge":         "brv",
		"copper_creek":          "ccv",
		"grand_californian":     "gcv",
		"grand_floridian":       "gfv",
		"hilton_head":           "hhr",
		"old_key_west":          "okw",
		"polynesian_villas":     "pvv",
		"riviera":               "riv",
		"saratoga_springs":      "ssr",
		"vero_beach":            "vbr",
	}
	ResortCodeToDVCRentalResortName = map[string]string{
		"akv": "animal_kingdom_villas",
		"aus": "aulani",
		"blt": "bay_lake_tower",
		"bcv": "beach_club_villas",
		"bwv": "boardwalk_villas",
		"brv": "boulder_ridge",
		"ccv": "copper_creek",
		"gcv": "grand_californian",
		"gfv": "grand_floridian",
		"hhr": "hilton_head",
		"okw": "old_key_west",
		"pvv": "polynesian_villas",
		"riv": "riviera",
		"ssr": "saratoga_springs",
		"vbr": "vero_beach",
	}
	DVCRentalRoomTypesToCode = map[string]struct {
		ViewRoomMap map[string]map[string]string
	}{
		"aul": {
			ViewRoomMap: map[string]map[string]string{
				"standard_view": {
					"hotel_room":                "htl",
					"deluxe_studio":             "dss",
					"one_bedroom_villa":         "1bs",
					"two_bedroom_villa":         "2bs",
					"three_bedroom_grand_villa": "3bs",
				},
				"island_gardens_view": {
					"deluxe_studio":             "dsi",
					"one_bedroom_villa":         "1bi",
					"two_bedroom_villa":         "2bi",
					"three_bedroom_grand_villa": "3bi",
				},
				"poolside_gardens_view": {
					"deluxe_studio":             "dsp",
					"one_bedroom_villa":         "1bp",
					"two_bedroom_villa":         "2bp",
					"three_bedroom_grand_villa": "3bp",
				},
				"ocean_view": {
					"deluxe_studio":             "dso",
					"one_bedroom_villa":         "1bo",
					"two_bedroom_villa":         "2bo",
					"three_bedroom_grand_villa": "3bo",
				},
			},
		},
		"ssr": {
			ViewRoomMap: map[string]map[string]string{
				"standard_view": {
					"deluxe_studio":             "dss",
					"one_bedroom_villa":         "1bs",
					"two_bedroom_villa":         "2bs",
					"three_bedroom_grand_villa": "3bs",
					"treehouse":                 "3bt",
				},
				"preferred_view": {
					"deluxe_studio":             "dsp",
					"one_bedroom_villa":         "1bp",
					"two_bedroom_villa":         "2bp",
					"three_bedroom_grand_villa": "3bp",
				},
			},
		},
		"blt": {
			ViewRoomMap: map[string]map[string]string{
				"standard_view": {
					"deluxe_studio":     "dss",
					"one_bedroom_villa": "1bs",
					"two_bedroom_villa": "2bs",
				},
				"lake_view": {
					"deluxe_studio":             "dsl",
					"one_bedroom_villa":         "1bl",
					"two_bedroom_villa":         "2bl",
					"three_bedroom_grand_villa": "3bl",
				},
				"theme_park_view": {
					"deluxe_studio":             "dsp",
					"one_bedroom_villa":         "1bp",
					"two_bedroom_villa":         "2bp",
					"three_bedroom_grand_villa": "3bp",
				},
			},
		},
	}
)

type SeasonYearViewTierDayTypeRoomPrices map[string]string
type SeasonYearViewTierPrices struct {
	SunThur SeasonYearViewTierDayTypeRoomPrices `json:"sunThur"`
	FriSat  SeasonYearViewTierDayTypeRoomPrices `json:"friSat"`
}
type SeasonYearViewPrices map[string]SeasonYearViewTierPrices
type SeasonYearPrices map[string]SeasonYearViewPrices
type SeasonPrices map[string]SeasonYearPrices

type View struct {
	View string `json:"view"`
}

type Resort struct {
	Id             string                  `json:"id"`
	State          string                  `json:"state"`
	Name           string                  `json:"name"`
	RoomTypes      map[string]bool         `json:"room_types"`
	Views          map[string]View         `json:"views"`
	Ordering       string                  `json:"ordering"`
	CheckedOut     string                  `json:"checked_out"`
	CheckedOutTime string                  `json:"checked_out_time"`
	CreatedBy      string                  `json:"created_by"`
	ModifiedBy     string                  `json:"modified_by"`
	Prices         map[string]SeasonPrices `json:"prices"`
}

type DateRangeList []*model.DateRange

func (drl DateRangeList) Len() int {
	return len(drl)
}

func (drl DateRangeList) Less(i, j int) bool {
	return drl[i].StartDate.Before(drl[j].StartDate)
}

func (drl DateRangeList) Swap(i, j int) {
	tmp := drl[i]
	drl[i] = drl[j]
	drl[j] = tmp
}

func (drl DateRangeList) String() string {
	buf := bytes.NewBuffer([]byte{})
	fmt.Fprint(buf, "[")
	for i, r := range drl {
		if i > 0 {
			fmt.Fprint(buf, ",")
		}
		fmt.Fprintf(buf, "%s", r)
	}
	fmt.Fprint(buf, "]")
	return buf.String()
}

func (drl DateRangeList) Equal(ddrl DateRangeList) bool {
	if len(drl) == len(ddrl) {
		for idx, dr := range drl {
			ddr := ddrl[idx]
			if !dr.Equal(ddr) {
				return false
			}
		}
		return true
	}
	return false
}

type TierRange []string

func (tr TierRange) GetRange() (dateRange model.DateRange, err error) {
	if tr[0] == "-" {
		return
	}
	dateRange.StartDate, err = model.ParseTierDate(tr[0])
	if err != nil {
		return
	}
	dateRange.EndDate, err = model.ParseTierDate(tr[1])
	if err != nil {
		return
	}
	return
}

type SeasonTier []TierRange

func (sy SeasonTier) GetRanges() (ranges DateRangeList, err error) {
	ranges = make(DateRangeList, 0, len(sy))
	var zeroTime time.Time
	for _, syr := range sy {
		var dateRange model.DateRange
		dateRange, err = syr.GetRange()
		if err != nil {
			break
		}
		if dateRange.StartDate != model.TierDate(zeroTime) {
			ranges = append(ranges, &dateRange)
		}
	}
	ranges = SortAndMergeDateRange(ranges)
	return
}

type SeasonYear map[string]SeasonTier
type Season map[string]SeasonYear
type Seasons map[string]Season

type Resorts map[string]Resort

func LoadSeasonsAndResorts(seasonsFileName, resortsFileName string) (seasons Seasons, resorts Resorts, err error) {
	var seasonsFile *os.File
	seasonsFile, err = os.Open(seasonsFileName)
	if err == nil {
		defer seasonsFile.Close()
		seasonsDecoder := json.NewDecoder(seasonsFile)
		seasons = make(Seasons)
		err = seasonsDecoder.Decode(&seasons)
		if err == nil {
			var resortsFile *os.File
			resortsFile, err = os.Open(resortsFileName)
			if err == nil {
				defer resortsFile.Close()
				resortsDecoder := json.NewDecoder(resortsFile)
				resorts = make(Resorts)
				err = resortsDecoder.Decode(&resorts)
			}
		}
	}
	return
}

func SortAndMergeDateRange(in DateRangeList) DateRangeList {
	sort.Sort(in)
	for i := 0; i < len(in); i++ {
		if i+1 < len(in) && in[i].EndDate.Add(time.Hour*24) == in[i+1].StartDate {
			in[i].EndDate = in[i+1].EndDate
			in = append(in[:i+1], in[i+2:]...)
		}
	}
	return in
}

func ConvertDVCRentalCharts(resortsToLoad []string, seasonsFileName, resortsFileName string) error {
	seasons, resorts, err := LoadSeasonsAndResorts(seasonsFileName, resortsFileName)
	if err != nil {
		return err
	}
	for DvcRResortName, DvcRResort := range resorts {
		resortCode := DVCRentalResortNameToCode[DvcRResortName]
		if !util.Contains(resortsToLoad, resortCode) {
			continue
		}
		for seasonId, seasonPrices := range DvcRResort.Prices {
			for yearStr, seasonYearPrices := range seasonPrices {
				year, _ := strconv.Atoi(yearStr)
				pointChart := &model.PointChart{
					Resort:    DvcRResortName,
					Code:      resortCode,
					Year:      year,
					RoomTypes: make([]*model.RoomType, 0, 20),
					Tiers:     make([]*model.Tier, 0, 10),
				}
				roomTypeMap := make(map[string]*model.RoomType)
				tierMap := make(map[string]*model.Tier)
				for tierName, DvcRTier := range seasons[seasonId][yearStr] {
					ranges, err := DvcRTier.GetRanges()
					if err != nil {
						return fmt.Errorf("error parsing ranges for seasonId = %s, year = %d, tierName = %s: %s",
							seasonId, year, tierName, err.Error())
					}
					tier := &model.Tier{
						DateRanges:     ranges,
						RoomTypePoints: make(map[string]*model.Points),
					}
					tierMap[tierName] = tier
				}
				for viewType, viewPrices := range seasonYearPrices {
					for tierName, tierPrices := range viewPrices {
						tier := tierMap[tierName]
						if tier == nil {
							continue
							//return fmt.Errorf("cannot find tier: %s - resortCode=%s, seasonId=%s, year=%d, viewType=%s", tierName, resortCode, seasonId, year, viewType)
						}
						for roomType, pointsStr := range tierPrices.SunThur {
							roomCode, ok := DVCRentalRoomTypesToCode[resortCode].ViewRoomMap[viewType][roomType]
							if !ok {
								continue
								//return fmt.Errorf("cannot get roomCode: resportCode=%s, viewType=%s, roomType=%s", resortCode, viewType, roomType)
							}
							rt := roomTypeMap[roomCode]
							if rt == nil {
								roomTypeMap[roomCode] = &model.RoomType{
									Name: roomType,
									Code: roomCode,
								}
							}
							points := tier.RoomTypePoints[roomCode]
							if points == nil {
								points = new(model.Points)
								tier.RoomTypePoints[roomCode] = points
							}
							points.Weekday, _ = strconv.Atoi(pointsStr)
						}
						for roomType, pointsStr := range tierPrices.FriSat {
							roomCode, ok := DVCRentalRoomTypesToCode[resortCode].ViewRoomMap[viewType][roomType]
							if !ok {
								continue
							}
							rt := roomTypeMap[roomCode]
							if rt == nil {
								roomTypeMap[roomCode] = &model.RoomType{
									Code: roomCode,
								}

							}
							points := tier.RoomTypePoints[roomCode]
							if points == nil {
								points = new(model.Points)
								tier.RoomTypePoints[roomCode] = points
							}
							points.Weekend, _ = strconv.Atoi(pointsStr)
						}
					}
				}
				for _, tier := range tierMap {
					pointChart.Tiers = append(pointChart.Tiers, tier)
				}
				for _, roomType := range roomTypeMap {
					pointChart.RoomTypes = append(pointChart.RoomTypes, roomType)
				}
				ctx := context.Background()
				err := db.SavePointChart(ctx, pointChart)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
