package api

import (
	"dvccalc/chart"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	ValidationTagStayLongerThan1Year = "stayLongerThan1Year"
)

func stayRangeValidator(sl validator.StructLevel) {
	stay, ok := sl.Current().Interface().(chart.Stay)
	if !ok {
		log.Printf("cant get stay struct")
		return
	}
	if stay.To.Sub(stay.From) > time.Hour*24*365 {
		sl.ReportError(stay.To, "to", "To", ValidationTagStayLongerThan1Year, stay.From.Format(time.RFC3339))
	}
}

func ReportErrors(err error) gin.H {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		report := make(gin.H)
		errors := make([]string, len(validationErrors))
		report["errors"] = errors
		for i, fieldError := range validationErrors {
			switch fieldError.Tag() {
			case ValidationTagStayLongerThan1Year:
				from, _ := time.Parse(time.RFC3339, fieldError.Param())
				to := fieldError.Value().(time.Time)
				errors[i] = fmt.Sprintf("stay cannot be longer than 1 year, %s - %s", from.Format("01-02-2006"), to.Format("01-02-2006"))
			case "ltefield":
				errors[i] = fmt.Sprintf("'from' date must be before 'to' date")
			case "len":
				errors[i] = fmt.Sprintf("%s must have length of %s got %+v", fieldError.Field(), fieldError.Param(), fieldError.Value())
			default:
				errors[i] = fmt.Sprintf("validation failed: field=%s, value=%+v, tag=%s, param=%s", fieldError.Field(), fieldError.Value(), fieldError.Tag(), fieldError.Param())
			}
		}
		report["msg"] = "validation failed"
		return report
	}
	return gin.H{"msg": err.Error()}
}

func GetStay(context *gin.Context) {
	var stay chart.Stay
	if err := context.ShouldBindUri(&stay); err != nil {
		context.JSON(http.StatusBadRequest, ReportErrors(err))
		return
	}
	if err := context.ShouldBindQuery(&stay); err != nil {
		context.JSON(http.StatusBadRequest, ReportErrors(err))
		return
	}
	log.Printf("GET /stay %+v", stay)
	result, err := chart.StayQuery(context, &stay)
	if err != nil {
		if strings.HasPrefix(err.Error(), chart.ErrorChartNotFound) {
			context.JSON(http.StatusNotFound, ReportErrors(err))
		} else {
			context.JSON(http.StatusInternalServerError, ReportErrors(err))
		}
		return
	}

	context.JSON(http.StatusOK, result)
}
