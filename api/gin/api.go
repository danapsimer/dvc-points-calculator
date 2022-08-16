package gin

import (
	"fmt"
	"github.com/danapsimer/dvc-points-calculator/db"
	"github.com/danapsimer/dvc-points-calculator/model"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"net/http"
)

// GET /stay/:from/:to?resort=<resort list>
// TODO: &sleepsMin=#&sleepsMax=#&bedroomsMin=#&bedroomsMax=#&roomClass=<class list>&pointsMax=#&sort=<field list>&sortDir=<asc|dsc>
// TODO: quality = preferred | standard
// TODO: sort = point | resort | type | bedrooms | sleeps
//{
//	"from": dateTime,
//	"to": dateTime,
//  "resort": [ string ]
//  "results": [
//    {
//      "resort": string,
//      "roomType": string,
//      "points": number
//    }
//  ]
//}

type ServiceConfig struct {
	ListenAddresses       []string
	GoogleProjectID       string
	GoogleCredentialsFile string
}

var (
	router        *gin.Engine
	defaultConfig = &ServiceConfig{
		[]string{":8080"},
		"dvc-points-calculator-qa",
		"./google-credentials.json",
	}
	config = defaultConfig
)

func notImplemented(context *gin.Context) {
	context.JSON(http.StatusNotImplemented, gin.H{"msg": fmt.Sprintf("%s not implemented", context.FullPath())})
}

func Start() error {
	router = gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterStructValidation(stayRangeValidator, model.Stay{})
	} else {
		panic("unable to register validators")
	}
	router.GET("/stay/:from/:to", GetStay)
	router.GET("/resort", GetResorts)
	router.POST("/resort", notImplemented)
	router.GET("/resort/:resortCode", notImplemented)
	router.PUT("/resort/:resortCode", UpdateResort)
	router.DELETE("/resort/:resortCode", notImplemented)
	router.GET("/resort/:resortCode/year/:year", GetResortYear)
	router.GET("/chart", notImplemented)
	router.GET("/chart/:resortCode", notImplemented)
	router.GET("/chart/:resortCode/:year", notImplemented)
	router.PUT("/chart/:resortCode/:year", notImplemented)
	router.DELETE("/chart/:resortCode/:year", notImplemented)
	err := db.InitDatastore(config.GoogleProjectID, config.GoogleCredentialsFile)
	if err != nil {
		return fmt.Errorf("error while initializing google datastore: %+v", err)
	}
	router.Run(config.ListenAddresses...)
	return nil
}
