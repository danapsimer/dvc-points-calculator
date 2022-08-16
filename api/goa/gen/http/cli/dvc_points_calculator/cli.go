// Code generated by goa v3.8.2, DO NOT EDIT.
//
// dvc-points-calculator HTTP client CLI support package
//
// Command:
// $ goa gen github.com/danapsimer/dvc-points-calculator/api/goa/design -o
// api/goa

package cli

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	dvcpointscalculatorc "github.com/danapsimer/dvc-points-calculator/api/goa/gen/http/dvc_points_calculator/client"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//	command (subcommand1|subcommand2|...)
func UsageCommands() string {
	return `dvc-points-calculator (get-resorts|get-resort)
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` dvc-points-calculator get-resorts` + "\n" +
		""
}

// ParseEndpoint returns the endpoint and payload as specified on the command
// line.
func ParseEndpoint(
	scheme, host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restore bool,
) (goa.Endpoint, interface{}, error) {
	var (
		dvcPointsCalculatorFlags = flag.NewFlagSet("dvc-points-calculator", flag.ContinueOnError)

		dvcPointsCalculatorGetResortsFlags = flag.NewFlagSet("get-resorts", flag.ExitOnError)

		dvcPointsCalculatorGetResortFlags          = flag.NewFlagSet("get-resort", flag.ExitOnError)
		dvcPointsCalculatorGetResortResortCodeFlag = dvcPointsCalculatorGetResortFlags.String("resort-code", "REQUIRED", "the resort's code")
	)
	dvcPointsCalculatorFlags.Usage = dvcPointsCalculatorUsage
	dvcPointsCalculatorGetResortsFlags.Usage = dvcPointsCalculatorGetResortsUsage
	dvcPointsCalculatorGetResortFlags.Usage = dvcPointsCalculatorGetResortUsage

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		return nil, nil, err
	}

	if flag.NArg() < 2 { // two non flag args are required: SERVICE and ENDPOINT (aka COMMAND)
		return nil, nil, fmt.Errorf("not enough arguments")
	}

	var (
		svcn string
		svcf *flag.FlagSet
	)
	{
		svcn = flag.Arg(0)
		switch svcn {
		case "dvc-points-calculator":
			svcf = dvcPointsCalculatorFlags
		default:
			return nil, nil, fmt.Errorf("unknown service %q", svcn)
		}
	}
	if err := svcf.Parse(flag.Args()[1:]); err != nil {
		return nil, nil, err
	}

	var (
		epn string
		epf *flag.FlagSet
	)
	{
		epn = svcf.Arg(0)
		switch svcn {
		case "dvc-points-calculator":
			switch epn {
			case "get-resorts":
				epf = dvcPointsCalculatorGetResortsFlags

			case "get-resort":
				epf = dvcPointsCalculatorGetResortFlags

			}

		}
	}
	if epf == nil {
		return nil, nil, fmt.Errorf("unknown %q endpoint %q", svcn, epn)
	}

	// Parse endpoint flags if any
	if svcf.NArg() > 1 {
		if err := epf.Parse(svcf.Args()[1:]); err != nil {
			return nil, nil, err
		}
	}

	var (
		data     interface{}
		endpoint goa.Endpoint
		err      error
	)
	{
		switch svcn {
		case "dvc-points-calculator":
			c := dvcpointscalculatorc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "get-resorts":
				endpoint = c.GetResorts()
				data = nil
			case "get-resort":
				endpoint = c.GetResort()
				data, err = dvcpointscalculatorc.BuildGetResortPayload(*dvcPointsCalculatorGetResortResortCodeFlag)
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// dvc-points-calculatorUsage displays the usage of the dvc-points-calculator
// command and its subcommands.
func dvcPointsCalculatorUsage() {
	fmt.Fprintf(os.Stderr, `provides resources for manipulating resorts, point charts, and querying stays
Usage:
    %[1]s [globalflags] dvc-points-calculator COMMAND [flags]

COMMAND:
    get-resorts: GetResorts implements GetResorts.
    get-resort: GetResort implements GetResort.

Additional help:
    %[1]s dvc-points-calculator COMMAND --help
`, os.Args[0])
}
func dvcPointsCalculatorGetResortsUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] dvc-points-calculator get-resorts

GetResorts implements GetResorts.

Example:
    %[1]s dvc-points-calculator get-resorts
`, os.Args[0])
}

func dvcPointsCalculatorGetResortUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] dvc-points-calculator get-resort -resort-code STRING

GetResort implements GetResort.
    -resort-code STRING: the resort's code

Example:
    %[1]s dvc-points-calculator get-resort --resort-code "ssr"
`, os.Args[0])
}
