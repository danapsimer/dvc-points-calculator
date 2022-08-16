package design

import (
	"github.com/danapsimer/dvc-points-calculator/model"
	. "goa.design/goa/v3/dsl"
)

var _ = API("Points", func() {
	Title("DVC Points Calculator")
	Description("Backend service for DVC Points Calculator UI")
	Server("points", func() {
		Host("localhost", func() {
			URI("http://localhost:8080")
		})
	})
})

func ChartYear() {
	Minimum(1990)
	Maximum(2100)
}

func RoomCode() {
	Pattern("[a-z0-9]{3}")
	Example("1bp")
}

func ResortCode() {
	Pattern("[a-z]{3}")
	Example("ssr")
}

func Sleeps() {
	Example(5)
	Minimum(1)
	Maximum(12)
}

func Bedrooms() {
	Example(1)
	Minimum(0)
	Maximum(3)
}

func Beds() {
	Example(2)
	Minimum(2)
	Maximum(6)
}

func RoomTypeName() {
	Example("One Bedroom Villa (Preferred View)")
}
func ResortName() {
	Example("Disney's Saratoga Springs Resort & Spa")
}

var RoomType = Type("RoomType", func() {
	Attribute("code", String, "room type's code", RoomCode)
	Attribute("name", String, "room type's name", RoomTypeName)
	Attribute("sleeps", Int, "max room capacity", Sleeps)
	Attribute("bedrooms", Int, "number of bedrooms", Bedrooms)
	Attribute("beds", Int, "number of beds", Beds)
	ConvertTo(model.RoomType{})
	CreateFrom(model.RoomType{})
})

var Resort = Type("Resort", func() {
	Attribute("code", String, "resort's code", ResortCode)
	Attribute("name", String, "resort's name", ResortName)
	Attribute("roomTypes", ArrayOf(RoomType))
	ConvertTo(model.Resort{})
})

var ResortResponse = ResultType("application/vnd.dvc.point.calculator.resort", "ResortResult", func() {
	Extend(Resort)
	View("default", func() {
		Attribute("code")
		Attribute("name")
		Attribute("roomTypes")
	})
	View("resortOnly", func() {
		Attribute("code")
		Attribute("name")
	})
	View("resortUpdate", func() {
		Attribute("name")
	})
	CreateFrom(model.Resort{})
})

func MonthDay() {
	Pattern("\\d{1,2}-\\d{1,2}")
}

var TierDateRange = Type("TierDateRange", func() {
	Attribute("startDate", String, "start date", MonthDay)
	Attribute("endDate", String, "end date", MonthDay)
})

var TierRoomTypePoints = Type("TierRoomTypePoints", func() {
	Attribute("weekday", Int, "points for Sunday - Thursday")
	Attribute("weekend", Int, "points for Friday - Saturday")
})

var Tier = Type("Tier", func() {
	Attribute("dateRanges", ArrayOf(TierDateRange))
	Attribute("roomTypePoints", MapOf(String, TierRoomTypePoints))
})

var PointChart = Type("PointChart", func() {
	Attribute("code", String, "resort's code", ResortCode)
	Attribute("resort", String, "resort's code", ResortName)
	Attribute("roomTypes", ArrayOf(RoomType))
	Attribute("tiers", ArrayOf(Tier))
})

var ResortYearResponse = ResultType("application/vnd.dvc.point.calculator.resortYear", "ResortYearResult", func() {
	Extend(ResortResponse)
	Attribute("year", Int, "the year the resort info is for.", ChartYear)
	View("default", func() {
		Attribute("code")
		Attribute("name")
		Attribute("year")
		Attribute("roomTypes")
	})
})

var ListResorts = CollectionOf(ResortResponse, func() {
	View("resortOnly")
})

var _ = Service("Points", func() {
	Description("provides resources for manipulating resorts, point charts, and querying stays")
	Method("GetResorts", func() {
		Result(ListResorts)
		HTTP(func() {
			GET("/resort")
			Response(StatusOK)
		})
	})
	Method("GetResort", func() {
		Result(ResortResponse)
		Payload(func() {
			Attribute("resortCode", String, "the resort's code", ResortCode)
			Required("resortCode")
		})
		Error("not_found")
		HTTP(func() {
			GET("/resort/{resortCode}")
			Response(StatusOK)
			Response("not_found", StatusNotFound)
		})
	})
	Method("GetResortYear", func() {
		Result(ResortYearResponse)
		Payload(func() {
			Attribute("resortCode", String, "the resort's code", ResortCode)
			Attribute("year", Int, "the year", ChartYear)
			Required("resortCode", "year")
		})
		Error("not_found")
		HTTP(func() {
			GET("/resort/{resortCode}/year/{year}")
			Response(StatusOK)
			Response("not_found", StatusNotFound)
		})
	})
	Method("GetPointChart", func() {
		Result(PointChart)
		Payload(func() {
			Attribute("resortCode", String, "the resort's code", ResortCode)
			Attribute("year", Int, "the year", ChartYear)
			Required("resortCode", "year")
		})
		Error("not_found")
		HTTP(func() {
			GET("chart/{resortCode}/{year}")
			Response(StatusOK)
			Response("not_found", StatusNotFound)
		})
	})
})
