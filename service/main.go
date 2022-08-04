package main

import (
	"dvccalc/api"
	"log"
)

func main() {
	err := api.Start()
	if err != nil {
		log.Fatalf("exiting with an error: %+v", err)
	}
}
