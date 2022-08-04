package api

import (
	"dvccalc/chart"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type resortQuery struct {
	ResortCode string `json:"resortCode" uri:"resortCode"`
	Year       int    `json:"year" uri:"year"`
}

type resortInfo struct {
	resortQuery
	ResortName string           `json:"resortName"`
	RoomTypes  []chart.RoomType `json:"roomTypes"`
}

func GetResort(context *gin.Context) {
	rq := &resortQuery{}
	if err := context.BindUri(rq); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err": err, "msg": "request is malformed"})
		return
	}
	pc, err := chart.LoadPointChartByCodeAndYear(context, rq.ResortCode, rq.Year)
	if err != nil {
		if strings.HasPrefix(err.Error(), chart.ErrorChartNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"err": err, "msg": "chart not found"})
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{"err": err, "msg": "unexpected error"})
		}
		return
	}
	result := &resortInfo{*rq, pc.Resort, pc.RoomTypes}
	context.JSON(http.StatusOK, result)
}
