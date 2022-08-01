package api

import "github.com/gin-gonic/gin"

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
	ListenAddresses []string
}

var (
	router        *gin.Engine
	defaultConfig = &ServiceConfig{[]string{":8080"}}
	config        = defaultConfig
)

func init() {
	router = gin.Default()
	router.GET("/stay/:from/:to", GetStay)
}

func Start() {
	router.Run(config.ListenAddresses...)
}
