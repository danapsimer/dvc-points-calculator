package design

import (
	"github.com/danapsimer/dvc-points-calculator/model"
	. "goa.design/goa/v3/dsl"
)

var _ = API("dvcPointsCalculator", func() {
	Title("DVC Points Calculator")
	Description("Backend service for DVC Points Calculator UI")
	Server("dvc-points-calculator", func() {
		Host("localhost", func() {
			URI("http://localhost:8080")
		})
	})
})

func RoomCode() {
	Pattern("[a-z0-9]{3}")
	Example("1bp")
}

func ResortCode() {
	Pattern("[a-z]{3}")
	Example("ssr")
}

var RoomType = Type("RoomType", func() {
	Attribute("code", String, "room type's code", RoomCode)
	Attribute("name", String, "room type's name", func() {
		Example("One Bedroom Villa (Preferred View)")
	})
	Attribute("sleeps", Int, "max room capacity", func() {
		Example(5)
	})
	Attribute("bedrooms", Int, "number of bedrooms", func() {
		Example(1)
	})
	Attribute("beds", Int, "number of beds", func() {
		Example(3)
	})
	ConvertTo(model.RoomType{})
	CreateFrom(model.RoomType{})
})

var Resort = Type("Resort", func() {
	Attribute("code", String, "resort's code", ResortCode)
	Attribute("name", String, "resort's name", func() {
		Example("Disney's Saratoga Springs Resort & Spa")
	})
	Attribute("roomTypes", ArrayOf(RoomType))
	ConvertTo(model.Resort{})
})

var ResortResponse = ResultType("application/vnd.dvcPointCalculator.resort", func() {
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

var ListResorts = CollectionOf(ResortResponse, func() {
	View("resortOnly")
})

var _ = Service("dvcPointsCalculator", func() {
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
})
