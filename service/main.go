package main

import (
	"github.com/danapsimer/dvc-points-calculator/api/gin"
	"log"
)

func main() {
	err := gin.Start()
	if err != nil {
		log.Fatalf("exiting with an error: %+v", err)
	}
}
