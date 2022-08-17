// Code generated by goa v3.8.2, DO NOT EDIT.
//
// Points HTTP client encoders and decoders
//
// Command:
// $ goa gen github.com/danapsimer/dvc-points-calculator/api/goa/design -o
// api/goa

package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	points "github.com/danapsimer/dvc-points-calculator/api/goa/gen/points"
	pointsviews "github.com/danapsimer/dvc-points-calculator/api/goa/gen/points/views"
	goahttp "goa.design/goa/v3/http"
)

// BuildGetResortsRequest instantiates a HTTP request object with method and
// path set to call the "Points" service "GetResorts" endpoint
func (c *Client) BuildGetResortsRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: GetResortsPointsPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("Points", "GetResorts", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeGetResortsResponse returns a decoder for responses returned by the
// Points GetResorts endpoint. restoreBody controls whether the response body
// should be restored after having been read.
func DecodeGetResortsResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body GetResortsResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("Points", "GetResorts", err)
			}
			p := NewGetResortsResortResultCollectionOK(body)
			view := resp.Header.Get("goa-view")
			vres := pointsviews.ResortResultCollection{Projected: p, View: view}
			if err = pointsviews.ValidateResortResultCollection(vres); err != nil {
				return nil, goahttp.ErrValidationError("Points", "GetResorts", err)
			}
			res := points.NewResortResultCollection(vres)
			return res, nil
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("Points", "GetResorts", resp.StatusCode, string(body))
		}
	}
}

// BuildGetResortRequest instantiates a HTTP request object with method and
// path set to call the "Points" service "GetResort" endpoint
func (c *Client) BuildGetResortRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		resortCode string
	)
	{
		p, ok := v.(*points.GetResortPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("Points", "GetResort", "*points.GetResortPayload", v)
		}
		resortCode = p.ResortCode
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: GetResortPointsPath(resortCode)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("Points", "GetResort", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeGetResortResponse returns a decoder for responses returned by the
// Points GetResort endpoint. restoreBody controls whether the response body
// should be restored after having been read.
// DecodeGetResortResponse may return the following errors:
//   - "not_found" (type *goa.ServiceError): http.StatusNotFound
//   - error: internal error
func DecodeGetResortResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body GetResortResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("Points", "GetResort", err)
			}
			p := NewGetResortResortResultOK(&body)
			view := resp.Header.Get("goa-view")
			vres := &pointsviews.ResortResult{Projected: p, View: view}
			if err = pointsviews.ValidateResortResult(vres); err != nil {
				return nil, goahttp.ErrValidationError("Points", "GetResort", err)
			}
			res := points.NewResortResult(vres)
			return res, nil
		case http.StatusNotFound:
			var (
				body GetResortNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("Points", "GetResort", err)
			}
			err = ValidateGetResortNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("Points", "GetResort", err)
			}
			return nil, NewGetResortNotFound(&body)
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("Points", "GetResort", resp.StatusCode, string(body))
		}
	}
}

// BuildGetResortYearRequest instantiates a HTTP request object with method and
// path set to call the "Points" service "GetResortYear" endpoint
func (c *Client) BuildGetResortYearRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		resortCode string
		year       int
	)
	{
		p, ok := v.(*points.GetResortYearPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("Points", "GetResortYear", "*points.GetResortYearPayload", v)
		}
		resortCode = p.ResortCode
		year = p.Year
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: GetResortYearPointsPath(resortCode, year)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("Points", "GetResortYear", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeGetResortYearResponse returns a decoder for responses returned by the
// Points GetResortYear endpoint. restoreBody controls whether the response
// body should be restored after having been read.
// DecodeGetResortYearResponse may return the following errors:
//   - "not_found" (type *goa.ServiceError): http.StatusNotFound
//   - error: internal error
func DecodeGetResortYearResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body GetResortYearResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("Points", "GetResortYear", err)
			}
			p := NewGetResortYearResortYearResultOK(&body)
			view := "default"
			vres := &pointsviews.ResortYearResult{Projected: p, View: view}
			if err = pointsviews.ValidateResortYearResult(vres); err != nil {
				return nil, goahttp.ErrValidationError("Points", "GetResortYear", err)
			}
			res := points.NewResortYearResult(vres)
			return res, nil
		case http.StatusNotFound:
			var (
				body GetResortYearNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("Points", "GetResortYear", err)
			}
			err = ValidateGetResortYearNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("Points", "GetResortYear", err)
			}
			return nil, NewGetResortYearNotFound(&body)
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("Points", "GetResortYear", resp.StatusCode, string(body))
		}
	}
}

// BuildGetPointChartRequest instantiates a HTTP request object with method and
// path set to call the "Points" service "GetPointChart" endpoint
func (c *Client) BuildGetPointChartRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		resortCode string
		year       int
	)
	{
		p, ok := v.(*points.GetPointChartPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("Points", "GetPointChart", "*points.GetPointChartPayload", v)
		}
		resortCode = p.ResortCode
		year = p.Year
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: GetPointChartPointsPath(resortCode, year)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("Points", "GetPointChart", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeGetPointChartResponse returns a decoder for responses returned by the
// Points GetPointChart endpoint. restoreBody controls whether the response
// body should be restored after having been read.
// DecodeGetPointChartResponse may return the following errors:
//   - "not_found" (type *goa.ServiceError): http.StatusNotFound
//   - error: internal error
func DecodeGetPointChartResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body GetPointChartResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("Points", "GetPointChart", err)
			}
			err = ValidateGetPointChartResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("Points", "GetPointChart", err)
			}
			res := NewGetPointChartPointChartOK(&body)
			return res, nil
		case http.StatusNotFound:
			var (
				body GetPointChartNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("Points", "GetPointChart", err)
			}
			err = ValidateGetPointChartNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("Points", "GetPointChart", err)
			}
			return nil, NewGetPointChartNotFound(&body)
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("Points", "GetPointChart", resp.StatusCode, string(body))
		}
	}
}

// BuildQueryStayRequest instantiates a HTTP request object with method and
// path set to call the "Points" service "QueryStay" endpoint
func (c *Client) BuildQueryStayRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		from string
		to   string
	)
	{
		p, ok := v.(*points.Stay)
		if !ok {
			return nil, goahttp.ErrInvalidType("Points", "QueryStay", "*points.Stay", v)
		}
		from = p.From
		to = p.To
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: QueryStayPointsPath(from, to)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("Points", "QueryStay", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeQueryStayRequest returns an encoder for requests sent to the Points
// QueryStay server.
func EncodeQueryStayRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*points.Stay)
		if !ok {
			return goahttp.ErrInvalidType("Points", "QueryStay", "*points.Stay", v)
		}
		values := req.URL.Query()
		for _, value := range p.IncludeResorts {
			values.Add("includeResorts", value)
		}
		for _, value := range p.ExcludeResorts {
			values.Add("excludeResorts", value)
		}
		values.Add("minSleeps", fmt.Sprintf("%v", p.MinSleeps))
		values.Add("maxSleeps", fmt.Sprintf("%v", p.MaxSleeps))
		values.Add("minBedrooms", fmt.Sprintf("%v", p.MinBedrooms))
		values.Add("maxBedrooms", fmt.Sprintf("%v", p.MaxBedrooms))
		values.Add("minBeds", fmt.Sprintf("%v", p.MinBeds))
		values.Add("maxBeds", fmt.Sprintf("%v", p.MaxBeds))
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// DecodeQueryStayResponse returns a decoder for responses returned by the
// Points QueryStay endpoint. restoreBody controls whether the response body
// should be restored after having been read.
// DecodeQueryStayResponse may return the following errors:
//   - "invalid_input" (type *goa.ServiceError): http.StatusBadRequest
//   - error: internal error
func DecodeQueryStayResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body QueryStayResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("Points", "QueryStay", err)
			}
			err = ValidateQueryStayResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("Points", "QueryStay", err)
			}
			res := NewQueryStayStayResultOK(&body)
			return res, nil
		case http.StatusBadRequest:
			var (
				body QueryStayInvalidInputResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("Points", "QueryStay", err)
			}
			err = ValidateQueryStayInvalidInputResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("Points", "QueryStay", err)
			}
			return nil, NewQueryStayInvalidInput(&body)
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("Points", "QueryStay", resp.StatusCode, string(body))
		}
	}
}

// unmarshalResortResultResponseToPointsviewsResortResultView builds a value of
// type *pointsviews.ResortResultView from a value of type
// *ResortResultResponse.
func unmarshalResortResultResponseToPointsviewsResortResultView(v *ResortResultResponse) *pointsviews.ResortResultView {
	res := &pointsviews.ResortResultView{
		Code: v.Code,
		Name: v.Name,
	}
	if v.RoomTypes != nil {
		res.RoomTypes = make([]*pointsviews.RoomTypeView, len(v.RoomTypes))
		for i, val := range v.RoomTypes {
			res.RoomTypes[i] = unmarshalRoomTypeResponseToPointsviewsRoomTypeView(val)
		}
	}

	return res
}

// unmarshalRoomTypeResponseToPointsviewsRoomTypeView builds a value of type
// *pointsviews.RoomTypeView from a value of type *RoomTypeResponse.
func unmarshalRoomTypeResponseToPointsviewsRoomTypeView(v *RoomTypeResponse) *pointsviews.RoomTypeView {
	if v == nil {
		return nil
	}
	res := &pointsviews.RoomTypeView{
		Code:     v.Code,
		Name:     v.Name,
		Sleeps:   v.Sleeps,
		Bedrooms: v.Bedrooms,
		Beds:     v.Beds,
	}

	return res
}

// unmarshalRoomTypeResponseBodyToPointsviewsRoomTypeView builds a value of
// type *pointsviews.RoomTypeView from a value of type *RoomTypeResponseBody.
func unmarshalRoomTypeResponseBodyToPointsviewsRoomTypeView(v *RoomTypeResponseBody) *pointsviews.RoomTypeView {
	if v == nil {
		return nil
	}
	res := &pointsviews.RoomTypeView{
		Code:     v.Code,
		Name:     v.Name,
		Sleeps:   v.Sleeps,
		Bedrooms: v.Bedrooms,
		Beds:     v.Beds,
	}

	return res
}

// unmarshalRoomTypeResponseBodyToPointsRoomType builds a value of type
// *points.RoomType from a value of type *RoomTypeResponseBody.
func unmarshalRoomTypeResponseBodyToPointsRoomType(v *RoomTypeResponseBody) *points.RoomType {
	res := &points.RoomType{
		Code:     *v.Code,
		Name:     *v.Name,
		Sleeps:   *v.Sleeps,
		Bedrooms: *v.Bedrooms,
		Beds:     *v.Beds,
	}

	return res
}

// unmarshalTierResponseBodyToPointsTier builds a value of type *points.Tier
// from a value of type *TierResponseBody.
func unmarshalTierResponseBodyToPointsTier(v *TierResponseBody) *points.Tier {
	res := &points.Tier{}
	res.DateRanges = make([]*points.TierDateRange, len(v.DateRanges))
	for i, val := range v.DateRanges {
		res.DateRanges[i] = unmarshalTierDateRangeResponseBodyToPointsTierDateRange(val)
	}
	res.RoomTypePoints = make(map[string]*points.TierRoomTypePoints, len(v.RoomTypePoints))
	for key, val := range v.RoomTypePoints {
		tk := key
		res.RoomTypePoints[tk] = unmarshalTierRoomTypePointsResponseBodyToPointsTierRoomTypePoints(val)
	}

	return res
}

// unmarshalTierDateRangeResponseBodyToPointsTierDateRange builds a value of
// type *points.TierDateRange from a value of type *TierDateRangeResponseBody.
func unmarshalTierDateRangeResponseBodyToPointsTierDateRange(v *TierDateRangeResponseBody) *points.TierDateRange {
	res := &points.TierDateRange{
		StartDate: *v.StartDate,
		EndDate:   *v.EndDate,
	}

	return res
}

// unmarshalTierRoomTypePointsResponseBodyToPointsTierRoomTypePoints builds a
// value of type *points.TierRoomTypePoints from a value of type
// *TierRoomTypePointsResponseBody.
func unmarshalTierRoomTypePointsResponseBodyToPointsTierRoomTypePoints(v *TierRoomTypePointsResponseBody) *points.TierRoomTypePoints {
	res := &points.TierRoomTypePoints{
		Weekday: *v.Weekday,
		Weekend: *v.Weekend,
	}

	return res
}