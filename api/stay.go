package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Room struct {
	Resort   string `json:"resort"`
	RoomType string `json:"roomType"`
	Points   int    `json:"points"`
}

type Stay struct {
	From  time.Time `json:"from" uri:"from" binding:"required"`
	To    time.Time `json:"to" uri:"to" binding:"required"`
	Rooms []Room    `json:"rooms"`
}

func GetStay(context *gin.Context) {
	var stay Stay
	if err := context.ShouldBindUri(&stay); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	StayQuery(stay.From, stay.To)
}
