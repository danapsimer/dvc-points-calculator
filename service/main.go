package main

import (
	"github.com/danapsimer/dvc-points-calculator/api"
	"log"
)

func main() {
	err := api.Start()
	if err != nil {
		log.Fatalf("exiting with an error: %+v", err)
	}
}
