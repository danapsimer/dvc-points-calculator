// Code generated by goa v3.8.2, DO NOT EDIT.
//
// HTTP request path constructors for the dvcPointsCalculator service.
//
// Command:
// $ goa gen github.com/danapsimer/dvc-points-calculator/api/goa/design -o
// api/goa

package client

import (
	"fmt"
)

// GetResortsDvcPointsCalculatorPath returns the URL path to the dvcPointsCalculator service GetResorts HTTP endpoint.
func GetResortsDvcPointsCalculatorPath() string {
	return "/resort"
}

// GetResortDvcPointsCalculatorPath returns the URL path to the dvcPointsCalculator service GetResort HTTP endpoint.
func GetResortDvcPointsCalculatorPath(resortCode string) string {
	return fmt.Sprintf("/resort/%v", resortCode)
}