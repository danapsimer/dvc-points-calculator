package goa

import (
	"github.com/danapsimer/dvc-points-calculator/api/goa/gen/http/points/server"
	"github.com/danapsimer/dvc-points-calculator/api/goa/gen/points"
	"github.com/danapsimer/dvc-points-calculator/db"
	"github.com/spf13/viper"
	goahttp "goa.design/goa/v3/http"
	"net/http"
)

type Points struct {
	points.Service
	// Implements the points.Service interface
	// The function definitions are is other files within this package.
}

func Start() error {
	err := db.InitDatastore(viper.GetString("google.projectId"), viper.GetString("google.credentialsFile"))
	if err != nil {
		return err
	}
	s := &Points{}
	endpoints := points.NewEndpoints(s)                   // Create endpoints
	mux := goahttp.NewMuxer()                             // Create HTTP muxer
	dec := goahttp.RequestDecoder                         // Set HTTP request decoder
	enc := goahttp.ResponseEncoder                        // Set HTTP response encoder
	svr := server.New(endpoints, mux, dec, enc, nil, nil) // Create Goa HTTP server
	server.Mount(mux, svr)                                // Mount Goa server on mux
	listenAddresses := viper.GetStringSlice("service.listenAddresses")
	if listenAddresses == nil || len(listenAddresses) == 0 {
		listenAddresses = []string{"localhost:8080"}
	}
	httpServer := &http.Server{ // Create Go HTTP server
		Addr:    listenAddresses[0], // Configure server address
		Handler: mux,                // Set request handler
	}
	if err = httpServer.ListenAndServe(); err != nil { // Start HTTP server
		return err
	}
	return nil
}
