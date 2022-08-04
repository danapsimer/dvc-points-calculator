package api

import (
	"dvccalc/chart"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type resortQuery struct {
	ResortCode string `json:"resortCode" uri:"resortCode" binding:"required"`
	Year       int    `json:"year" uri:"year" binding:"required,len=4"`
}

type resortInfo struct {
	resortQuery
	ResortName string           `json:"resortName"`
	RoomTypes  []chart.RoomType `json:"roomTypes"`
}

func GetResort(context *gin.Context) {
	rq := &resortQuery{}
	if err := context.BindUri(rq); err != nil {
		context.JSON(http.StatusBadRequest, ReportErrors(err))
		return
	}
	pc, err := chart.LoadPointChartByCodeAndYear(context, rq.ResortCode, rq.Year)
	if err != nil {
		if strings.HasPrefix(err.Error(), chart.ErrorChartNotFound) {
			context.JSON(http.StatusNotFound, ReportErrors(err))
		} else {
			context.JSON(http.StatusInternalServerError, ReportErrors(err))
		}
		return
	}
	result := &resortInfo{*rq, pc.Resort, pc.RoomTypes}
	context.JSON(http.StatusOK, result)
}
