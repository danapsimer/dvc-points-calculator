// Code generated by goa v3.8.2, DO NOT EDIT.
//
// dvcPointsCalculator service
//
// Command:
// $ goa gen github.com/danapsimer/dvc-points-calculator/api/goa/design -o
// api/goa

package dvcpointscalculator

import (
	"context"

	dvcpointscalculatorviews "github.com/danapsimer/dvc-points-calculator/api/goa/gen/dvc_points_calculator/views"
	goa "goa.design/goa/v3/pkg"
)

// provides resources for manipulating resorts, point charts, and querying stays
type Service interface {
	// GetResorts implements GetResorts.
	// The "view" return value must have one of the following views
	//	- "default"
	//	- "resortOnly"
	//	- "resortUpdate"
	GetResorts(context.Context) (res DvcpointcalculatorResortCollection, view string, err error)
	// GetResort implements GetResort.
	// The "view" return value must have one of the following views
	//	- "default"
	//	- "resortOnly"
	//	- "resortUpdate"
	GetResort(context.Context, *GetResortPayload) (res *DvcpointcalculatorResort, view string, err error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "dvcPointsCalculator"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [2]string{"GetResorts", "GetResort"}

// DvcpointcalculatorResort is the result type of the dvcPointsCalculator
// service GetResort method.
type DvcpointcalculatorResort struct {
	// resort's code
	Code *string
	// resort's name
	Name      *string
	RoomTypes []*RoomType
}

// DvcpointcalculatorResortCollection is the result type of the
// dvcPointsCalculator service GetResorts method.
type DvcpointcalculatorResortCollection []*DvcpointcalculatorResort

// GetResortPayload is the payload type of the dvcPointsCalculator service
// GetResort method.
type GetResortPayload struct {
	// the resort's code
	ResortCode string
}

type RoomType struct {
	// room type's code
	Code *string
	// room type's name
	Name *string
	// max room capacity
	Sleeps *int
	// number of bedrooms
	Bedrooms *int
	// number of beds
	Beds *int
}

// MakeNotFound builds a goa.ServiceError from an error.
func MakeNotFound(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "not_found", false, false, false)
}

// NewDvcpointcalculatorResortCollection initializes result type
// DvcpointcalculatorResortCollection from viewed result type
// DvcpointcalculatorResortCollection.
func NewDvcpointcalculatorResortCollection(vres dvcpointscalculatorviews.DvcpointcalculatorResortCollection) DvcpointcalculatorResortCollection {
	var res DvcpointcalculatorResortCollection
	switch vres.View {
	case "default", "":
		res = newDvcpointcalculatorResortCollection(vres.Projected)
	case "resortOnly":
		res = newDvcpointcalculatorResortCollectionResortOnly(vres.Projected)
	case "resortUpdate":
		res = newDvcpointcalculatorResortCollectionResortUpdate(vres.Projected)
	}
	return res
}

// NewViewedDvcpointcalculatorResortCollection initializes viewed result type
// DvcpointcalculatorResortCollection from result type
// DvcpointcalculatorResortCollection using the given view.
func NewViewedDvcpointcalculatorResortCollection(res DvcpointcalculatorResortCollection, view string) dvcpointscalculatorviews.DvcpointcalculatorResortCollection {
	var vres dvcpointscalculatorviews.DvcpointcalculatorResortCollection
	switch view {
	case "default", "":
		p := newDvcpointcalculatorResortCollectionView(res)
		vres = dvcpointscalculatorviews.DvcpointcalculatorResortCollection{Projected: p, View: "default"}
	case "resortOnly":
		p := newDvcpointcalculatorResortCollectionViewResortOnly(res)
		vres = dvcpointscalculatorviews.DvcpointcalculatorResortCollection{Projected: p, View: "resortOnly"}
	case "resortUpdate":
		p := newDvcpointcalculatorResortCollectionViewResortUpdate(res)
		vres = dvcpointscalculatorviews.DvcpointcalculatorResortCollection{Projected: p, View: "resortUpdate"}
	}
	return vres
}

// NewDvcpointcalculatorResort initializes result type DvcpointcalculatorResort
// from viewed result type DvcpointcalculatorResort.
func NewDvcpointcalculatorResort(vres *dvcpointscalculatorviews.DvcpointcalculatorResort) *DvcpointcalculatorResort {
	var res *DvcpointcalculatorResort
	switch vres.View {
	case "default", "":
		res = newDvcpointcalculatorResort(vres.Projected)
	case "resortOnly":
		res = newDvcpointcalculatorResortResortOnly(vres.Projected)
	case "resortUpdate":
		res = newDvcpointcalculatorResortResortUpdate(vres.Projected)
	}
	return res
}

// NewViewedDvcpointcalculatorResort initializes viewed result type
// DvcpointcalculatorResort from result type DvcpointcalculatorResort using the
// given view.
func NewViewedDvcpointcalculatorResort(res *DvcpointcalculatorResort, view string) *dvcpointscalculatorviews.DvcpointcalculatorResort {
	var vres *dvcpointscalculatorviews.DvcpointcalculatorResort
	switch view {
	case "default", "":
		p := newDvcpointcalculatorResortView(res)
		vres = &dvcpointscalculatorviews.DvcpointcalculatorResort{Projected: p, View: "default"}
	case "resortOnly":
		p := newDvcpointcalculatorResortViewResortOnly(res)
		vres = &dvcpointscalculatorviews.DvcpointcalculatorResort{Projected: p, View: "resortOnly"}
	case "resortUpdate":
		p := newDvcpointcalculatorResortViewResortUpdate(res)
		vres = &dvcpointscalculatorviews.DvcpointcalculatorResort{Projected: p, View: "resortUpdate"}
	}
	return vres
}

// newDvcpointcalculatorResortCollection converts projected type
// DvcpointcalculatorResortCollection to service type
// DvcpointcalculatorResortCollection.
func newDvcpointcalculatorResortCollection(vres dvcpointscalculatorviews.DvcpointcalculatorResortCollectionView) DvcpointcalculatorResortCollection {
	res := make(DvcpointcalculatorResortCollection, len(vres))
	for i, n := range vres {
		res[i] = newDvcpointcalculatorResort(n)
	}
	return res
}

// newDvcpointcalculatorResortCollectionResortOnly converts projected type
// DvcpointcalculatorResortCollection to service type
// DvcpointcalculatorResortCollection.
func newDvcpointcalculatorResortCollectionResortOnly(vres dvcpointscalculatorviews.DvcpointcalculatorResortCollectionView) DvcpointcalculatorResortCollection {
	res := make(DvcpointcalculatorResortCollection, len(vres))
	for i, n := range vres {
		res[i] = newDvcpointcalculatorResortResortOnly(n)
	}
	return res
}

// newDvcpointcalculatorResortCollectionResortUpdate converts projected type
// DvcpointcalculatorResortCollection to service type
// DvcpointcalculatorResortCollection.
func newDvcpointcalculatorResortCollectionResortUpdate(vres dvcpointscalculatorviews.DvcpointcalculatorResortCollectionView) DvcpointcalculatorResortCollection {
	res := make(DvcpointcalculatorResortCollection, len(vres))
	for i, n := range vres {
		res[i] = newDvcpointcalculatorResortResortUpdate(n)
	}
	return res
}

// newDvcpointcalculatorResortCollectionView projects result type
// DvcpointcalculatorResortCollection to projected type
// DvcpointcalculatorResortCollectionView using the "default" view.
func newDvcpointcalculatorResortCollectionView(res DvcpointcalculatorResortCollection) dvcpointscalculatorviews.DvcpointcalculatorResortCollectionView {
	vres := make(dvcpointscalculatorviews.DvcpointcalculatorResortCollectionView, len(res))
	for i, n := range res {
		vres[i] = newDvcpointcalculatorResortView(n)
	}
	return vres
}

// newDvcpointcalculatorResortCollectionViewResortOnly projects result type
// DvcpointcalculatorResortCollection to projected type
// DvcpointcalculatorResortCollectionView using the "resortOnly" view.
func newDvcpointcalculatorResortCollectionViewResortOnly(res DvcpointcalculatorResortCollection) dvcpointscalculatorviews.DvcpointcalculatorResortCollectionView {
	vres := make(dvcpointscalculatorviews.DvcpointcalculatorResortCollectionView, len(res))
	for i, n := range res {
		vres[i] = newDvcpointcalculatorResortViewResortOnly(n)
	}
	return vres
}

// newDvcpointcalculatorResortCollectionViewResortUpdate projects result type
// DvcpointcalculatorResortCollection to projected type
// DvcpointcalculatorResortCollectionView using the "resortUpdate" view.
func newDvcpointcalculatorResortCollectionViewResortUpdate(res DvcpointcalculatorResortCollection) dvcpointscalculatorviews.DvcpointcalculatorResortCollectionView {
	vres := make(dvcpointscalculatorviews.DvcpointcalculatorResortCollectionView, len(res))
	for i, n := range res {
		vres[i] = newDvcpointcalculatorResortViewResortUpdate(n)
	}
	return vres
}

// newDvcpointcalculatorResort converts projected type DvcpointcalculatorResort
// to service type DvcpointcalculatorResort.
func newDvcpointcalculatorResort(vres *dvcpointscalculatorviews.DvcpointcalculatorResortView) *DvcpointcalculatorResort {
	res := &DvcpointcalculatorResort{
		Code: vres.Code,
		Name: vres.Name,
	}
	if vres.RoomTypes != nil {
		res.RoomTypes = make([]*RoomType, len(vres.RoomTypes))
		for i, val := range vres.RoomTypes {
			res.RoomTypes[i] = transformDvcpointscalculatorviewsRoomTypeViewToRoomType(val)
		}
	}
	return res
}

// newDvcpointcalculatorResortResortOnly converts projected type
// DvcpointcalculatorResort to service type DvcpointcalculatorResort.
func newDvcpointcalculatorResortResortOnly(vres *dvcpointscalculatorviews.DvcpointcalculatorResortView) *DvcpointcalculatorResort {
	res := &DvcpointcalculatorResort{
		Code: vres.Code,
		Name: vres.Name,
	}
	return res
}

// newDvcpointcalculatorResortResortUpdate converts projected type
// DvcpointcalculatorResort to service type DvcpointcalculatorResort.
func newDvcpointcalculatorResortResortUpdate(vres *dvcpointscalculatorviews.DvcpointcalculatorResortView) *DvcpointcalculatorResort {
	res := &DvcpointcalculatorResort{
		Name: vres.Name,
	}
	return res
}

// newDvcpointcalculatorResortView projects result type
// DvcpointcalculatorResort to projected type DvcpointcalculatorResortView
// using the "default" view.
func newDvcpointcalculatorResortView(res *DvcpointcalculatorResort) *dvcpointscalculatorviews.DvcpointcalculatorResortView {
	vres := &dvcpointscalculatorviews.DvcpointcalculatorResortView{
		Code: res.Code,
		Name: res.Name,
	}
	if res.RoomTypes != nil {
		vres.RoomTypes = make([]*dvcpointscalculatorviews.RoomTypeView, len(res.RoomTypes))
		for i, val := range res.RoomTypes {
			vres.RoomTypes[i] = transformRoomTypeToDvcpointscalculatorviewsRoomTypeView(val)
		}
	}
	return vres
}

// newDvcpointcalculatorResortViewResortOnly projects result type
// DvcpointcalculatorResort to projected type DvcpointcalculatorResortView
// using the "resortOnly" view.
func newDvcpointcalculatorResortViewResortOnly(res *DvcpointcalculatorResort) *dvcpointscalculatorviews.DvcpointcalculatorResortView {
	vres := &dvcpointscalculatorviews.DvcpointcalculatorResortView{
		Code: res.Code,
		Name: res.Name,
	}
	return vres
}

// newDvcpointcalculatorResortViewResortUpdate projects result type
// DvcpointcalculatorResort to projected type DvcpointcalculatorResortView
// using the "resortUpdate" view.
func newDvcpointcalculatorResortViewResortUpdate(res *DvcpointcalculatorResort) *dvcpointscalculatorviews.DvcpointcalculatorResortView {
	vres := &dvcpointscalculatorviews.DvcpointcalculatorResortView{
		Name: res.Name,
	}
	return vres
}

// transformDvcpointscalculatorviewsRoomTypeViewToRoomType builds a value of
// type *RoomType from a value of type *dvcpointscalculatorviews.RoomTypeView.
func transformDvcpointscalculatorviewsRoomTypeViewToRoomType(v *dvcpointscalculatorviews.RoomTypeView) *RoomType {
	if v == nil {
		return nil
	}
	res := &RoomType{
		Code:     v.Code,
		Name:     v.Name,
		Sleeps:   v.Sleeps,
		Bedrooms: v.Bedrooms,
		Beds:     v.Beds,
	}

	return res
}

// transformRoomTypeToDvcpointscalculatorviewsRoomTypeView builds a value of
// type *dvcpointscalculatorviews.RoomTypeView from a value of type *RoomType.
func transformRoomTypeToDvcpointscalculatorviewsRoomTypeView(v *RoomType) *dvcpointscalculatorviews.RoomTypeView {
	if v == nil {
		return nil
	}
	res := &dvcpointscalculatorviews.RoomTypeView{
		Code:     v.Code,
		Name:     v.Name,
		Sleeps:   v.Sleeps,
		Bedrooms: v.Bedrooms,
		Beds:     v.Beds,
	}

	return res
}