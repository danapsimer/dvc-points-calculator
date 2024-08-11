package design

import (
	"github.com/danapsimer/dvc-points-calculator/model"
	. "goa.design/goa/v3/dsl"
	"net/http"
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
	Required("code", "name", "sleeps", "bedrooms", "beds")
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
		Required("code", "name", "roomTypes")
	})
	View("resortOnly", func() {
		Attribute("code")
		Attribute("name")
		Required("code", "name")
	})
	View("resortUpdate", func() {
		Attribute("name")
		Required("name")
	})
	CreateFrom(model.Resort{})
})

func MonthDay() {
	Pattern("(0?1|0?2|0?3|0?4|0?5|0?6|0?7|0?8|0?9|10|11|12)-(0?1|0?2|0?3|0?4|0?5|0?6|0?7|0?8|0?9|10|11|12|13|14|15|16|17|18|19|20|21|22|23|24|25|26|27|28|29|30|31)")
	Meta("struct:field:type", "model.TierDate", "github.com/danapsimer/dvc-points-calculator/model")
	Example("1-15")
}

var TierDateRange = Type("TierDateRange", func() {
	Attribute("startDate", String, "start date", MonthDay)
	Attribute("endDate", String, "end date", MonthDay)
	Required("startDate", "endDate")
	ConvertTo(model.DateRange{})
	CreateFrom(model.DateRange{})
})

var TierRoomTypePoints = Type("TierRoomTypePoints", func() {
	Attribute("weekday", Int, "points for Sunday - Thursday", func() {
		Example(12)
	})
	Attribute("weekend", Int, "points for Friday - Saturday", func() {
		Example(17)
	})
	Required("weekday", "weekend")
	ConvertTo(model.Points{})
	CreateFrom(model.Points{})
})

var Tier = Type("Tier", func() {
	Attribute("dateRanges", ArrayOf(TierDateRange))
	Attribute("roomTypePoints", MapOf(String, TierRoomTypePoints), func() {
		Key(RoomCode)
	})
	Required("dateRanges", "roomTypePoints")
	ConvertTo(model.Tier{})
	CreateFrom(model.Tier{})
})

var PointChart = Type("PointChart", func() {
	Attribute("code", String, "resort's code", ResortCode)
	Attribute("resort", String, "resort's code", ResortName)
	Attribute("roomTypes", ArrayOf(RoomType))
	Attribute("tiers", ArrayOf(Tier))
	Required("code", "resort", "roomTypes", "tiers")
	ConvertTo(model.PointChart{})
	CreateFrom(model.PointChart{})
})

var ResortYearResponse = ResultType("application/vnd.dvc.point.calculator.resortYear", "ResortYearResult", func() {
	Extend(ResortResponse)
	Attribute("year", Int, "the year the resort info is for.", ChartYear)
	View("default", func() {
		Attribute("code")
		Attribute("name")
		Attribute("year")
		Attribute("roomTypes")
		Required("code", "name", "year", "roomTypes")
	})
})

var ListResorts = CollectionOf(ResortResponse, func() {
	View("resortOnly")
})

func StayDate() {
	Format(FormatDate)
	Example("2022-01-15")
}

func ResortCodeList() {
	Example([]string{"blt", "ssr", "akv"})
}

var Stay = Type("Stay", func() {
	Attribute("from", String, "Check-in Date", func() {
		StayDate()
		Example("2022-05-05")
	})
	Attribute("to", String, "Check-in Date", func() {
		StayDate()
		Example("2022-05-12")
	})
	Required("from", "to")
	Attribute("includeResorts", ArrayOf(String, ResortCode), "resorts to include in the search", ResortCodeList)
	Attribute("excludeResorts", ArrayOf(String, ResortCode), "resorts to exclude from the search", ResortCodeList)
	Attribute("minSleeps", Int, "the minimum capacity of room types to include", func() {
		Sleeps()
		Default(1)
	})
	Attribute("maxSleeps", Int, "the maximum capacity of room types to include", func() {
		Sleeps()
		Default(12)
	})
	Attribute("minBedrooms", Int, "the minimum number of bedrooms of room types to include", func() {
		Bedrooms()
		Default(0)
	})
	Attribute("maxBedrooms", Int, "the maximum number of bedrooms of room types to include", func() {
		Bedrooms()
		Default(3)
	})
	Attribute("minBeds", Int, "the minimum number of beds of room types to include", func() {
		Beds()
		Default(2)
	})
	Attribute("maxBeds", Int, "the maximum number of beds of room types to include", func() {
		Beds()
		Default(6)
	})
})

var StayResult = Type("StayResult", func() {
	Extend(Stay)
	Attribute("Rooms", MapOf(String, MapOf(String, Int)), func() {
		Example(map[string]map[string]int{
			"blt": {
				"dss": 12,
				"dsp": 15,
				"1bs": 20,
				"1bp": 25,
				"2bs": 30,
				"2bp": 35,
			},
		})
	})
	Required("Rooms")
})

var _ = Service("Points", func() {
	Description("provides resources for manipulating resorts, point charts, and querying stays")
	Error("not_found")
	Method("GetResorts", func() {
		Result(ListResorts)
		HTTP(func() {
			GET("/resort")
			Response(StatusOK)
		})
	})
	Method("GetResort", func() {
		Result(ResortResponse, func() {
			View("resortOnly")
		})
		Payload(func() {
			Attribute("resortCode", String, "the resort's code", ResortCode)
			Required("resortCode")
		})
		HTTP(func() {
			GET("/resort/{resortCode}")
			Response(StatusOK)
			Response("not_found", StatusNotFound)
		})
	})
	Method("PutResort", func() {
		Result(ResortResponse, func() {
			View("resortOnly")
		})
		Payload(func() {
			Attribute("resortCode", String, "the resort's code", ResortCode)
			Attribute("name", String, "The resort's name", ResortName)
			Required("resortCode", "name")
		})
		HTTP(func() {
			PUT("/resort/{resortCode}")
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
			GET("/chart/{resortCode}/{year}")
			Response(StatusOK)
			Response("not_found", StatusNotFound)
		})
	})
	Method("QueryStay", func() {
		Result(StayResult)
		Payload(Stay)
		Error("invalid_input")
		HTTP(func() {
			GET("/stay/{from}/{to}")
			Params(func() {
				Param("includeResorts")
				Param("excludeResorts")
				Param("minSleeps")
				Param("maxSleeps")
				Param("minBedrooms")
				Param("maxBedrooms")
				Param("minBeds")
				Param("maxBeds")
			})
			Response(StatusOK)
			Response("invalid_input", http.StatusBadRequest)
		})
	})
})
