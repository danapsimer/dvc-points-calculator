package api

import (
	"github.com/danapsimer/dvc-points-calculator/chart"
	"github.com/danapsimer/dvc-points-calculator/db"
	"github.com/danapsimer/dvc-points-calculator/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type resortQuery struct {
	ResortCode string `json:"resortCode" uri:"resortCode" binding:"required"`
	Year       string `json:"year" uri:"year" binding:"required,number,len=4"`
}

type resortInfo struct {
	resortQuery
	ResortName string            `json:"resortName"`
	RoomTypes  []*model.RoomType `json:"roomTypes"`
}

func GetResort(context *gin.Context) {
	rq := &resortQuery{}
	if err := context.BindUri(rq); err != nil {
		context.JSON(http.StatusBadRequest, ReportErrors(err))
		return
	}
	year, _ := strconv.Atoi(rq.Year)
	pc, err := db.GetPointChart(context, rq.ResortCode, year)
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
