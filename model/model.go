package model

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type TierDate time.Time

func ParseTierDate(dateStr string) (TierDate, error) {
	var tierDate TierDate

	parts := strings.Split(dateStr, "-")
	if len(parts) != 2 {
		return tierDate, fmt.Errorf("malformed tier date: %s - expected a month and day seperated by '-'", dateStr)
	}
	month, err := strconv.Atoi(parts[0])
	if err != nil {
		return tierDate, fmt.Errorf("malformed tier date: %s - expected a numeric month component: %s", dateStr, parts[0])
	}
	day, err := strconv.Atoi(parts[1])
	if err != nil {
		return tierDate, fmt.Errorf("malformed tier date: %s - expected a numeric day component: %s", dateStr, parts[1])
	}
	tierDate = TierDate(time.Date(time.Now().Year(), time.Month(month), day, 0, 0, 0, 0, time.UTC))
	return tierDate, nil
}

func (t TierDate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%02d-%02d\"", t.Month(), t.Day())), nil
}

func (t *TierDate) UnmarshalJSON(date []byte) error {
	dateStr := string(date)
	if !strings.HasPrefix(dateStr, "\"") || !strings.HasSuffix(dateStr, "\"") {
		return fmt.Errorf("malformed json string: %s", dateStr)
	}
	dateStr = dateStr[1 : len(dateStr)-1]
	d, err := ParseTierDate(dateStr)
	if err != nil {
		return err
	}
	*t = d
	return nil
}

func (t TierDate) Month() time.Month {
	return time.Time(t).Month()
}

func (t TierDate) Day() int {
	return time.Time(t).Day()
}

type DateRange struct {
	StartDate TierDate `json:"startDate"`
	EndDate   TierDate `json:"endDate"`
}

// Contains date must be truncated to midnight in UTC.
func (dr *DateRange) Contains(date time.Time) bool {
	from := time.Date(date.Year(), dr.StartDate.Month(), dr.StartDate.Day(), 0, 0, 0, 0, time.UTC)
	to := time.Date(date.Year(), dr.EndDate.Month(), dr.EndDate.Day(), 0, 0, 0, 0, time.UTC)
	return from.Equal(date) || (from.Before(date) && (to.Equal(date) || to.After(date)))
}

type RoomType struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Sleeps   int    `json:"sleeps"`
	Bedrooms int    `json:"bedrooms"`
	Beds     int    `json:"beds"`
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
	Code      string
	Year      int
	RoomTypes []RoomType
	Tiers     []Tier
}

type Stay struct {
	From           time.Time `json:"from" uri:"from" binding:"required,ltefield=To" time_format:"2006-01-02" time_utc:"1"`
	To             time.Time `json:"to" uri:"to" binding:"required" time_format:"2006-01-02" time_utc:"1"`
	IncludeResorts []string  `json:"includeResorts" form:"incResort"`
	ExcludeResorts []string  `json:"excludeResorts" form:"exResort"`
	MinSleeps      int       `json:"minSleeps" form:"minSleeps,default=1"`
	MaxSleeps      int       `json:"maxSleeps" form:"maxSleeps,default=12"`
	MinBedrooms    int       `json:"minBedrooms" form:"minBedrooms,default=0"`
	MaxBedrooms    int       `json:"maxBedrooms" form:"maxBedrooms,default=3"`
	MinBeds        int       `json:"minBeds" form:"minBeds,default=1"`
	MaxBeds        int       `json:"maxBeds" form:"maxBeds,default=6"`
}

type StayResult struct {
	Stay
	Rooms map[string]map[string]int `json:"rooms"`
}

func ReadPointChart(in io.Reader) (*PointChart, error) {
	decoder := json.NewDecoder(in)
	pc := &PointChart{}
	err := decoder.Decode(pc)
	if err != nil {
		return nil, fmt.Errorf("error decoding point chart: %s", err)
	}
	return pc, nil
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
	roomTypes := make([]string, 0, len(pc.RoomTypes))
	for _, roomType := range pc.RoomTypes {
		if stay.MinSleeps <= roomType.Sleeps && roomType.Sleeps <= stay.MaxSleeps &&
			stay.MinBedrooms <= roomType.Bedrooms && roomType.Bedrooms <= stay.MaxBedrooms &&
			stay.MinBeds <= roomType.Beds && roomType.Beds <= stay.MaxBeds {
			roomTypes = append(roomTypes, roomType.Code)
		}
	}
	for date := stay.From; date.Before(stay.To); date = date.Add(time.Hour * 24) {
		dayPoints, err := pc.GetPointsForDay(date, roomTypes...)
		if err != nil {
			return nil, fmt.Errorf("error getting points for a date: %+v - %s : %+v", stay, date.String(), err)
		}
		for rt, pts := range dayPoints {
			points[rt] += pts
		}
	}
	return points, nil
}

func (sr *StayResult) addPoints(resort, roomType string, points int) {
	if sr.Rooms[resort] == nil {
		sr.Rooms[resort] = make(map[string]int)
	}
	sr.Rooms[resort][roomType] += points
}

func (sr *StayResult) MergeResults(results ...*StayResult) {
	for _, result := range results {
		for resort, resortPoints := range result.Rooms {
			for roomType, roomTypePoints := range resortPoints {
				sr.addPoints(resort, roomType, roomTypePoints)
			}
		}
	}
}
