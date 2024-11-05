package health_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/GeoNet/kit/health"
)

var (
	healthCheckAged    = 5 * time.Second  //need to have a good heartbeat within this time
	healthCheckStartup = 5 * time.Second  //ignore heartbeat messages for this time after starting
	healthCheckTimeout = 30 * time.Second //health check timeout
	healthCheckService = ":7777"          //end point to listen to for SOH checks
	healthCheckPath    = "/soh"
)

func TestExistingSoh(t *testing.T) {
	checkPath := "https://api.geonet.org.nz/soh"
	if err := healthCheck(checkPath); err != nil {
		t.Error("should pass health check at start ")
	}
}

func TestHealth(t *testing.T) {
	checkPath := healthCheckService + healthCheckPath
	//1. should fail at start
	if err := healthCheck(checkPath); err == nil {
		t.Error("should fail health check at start ")
	}
	//2. start the process
	health := health.New(healthCheckService, healthCheckAged, healthCheckStartup)
	health.Ok()
	time.Sleep(1 * time.Millisecond) //let the service to start
	if err := healthCheck(checkPath); err != nil {
		t.Error("should pass health check after started ")
	}
	//3. test after healthCheckAged
	time.Sleep(healthCheckAged) //wait for the healthCheckAged
	if err := healthCheck(checkPath); err == nil {
		t.Errorf("should fail health check after %v", healthCheckAged)
	}

	//4. test after heartbeat
	health.Ok()
	if err := healthCheck(checkPath); err != nil {
		t.Error("should pass health check after heartbeat ")
	}
}

func TestHealthWithoutAgeCheck(t *testing.T) {
	healthCheckAged = 0 * time.Second
	healthCheckService = ":7778"
	checkPath := healthCheckService + healthCheckPath
	//1. start the process
	health := health.New(healthCheckService, healthCheckAged, healthCheckStartup)
	health.Ok()
	time.Sleep(1 * time.Millisecond) //let the service to start
	if err := healthCheck(checkPath); err != nil {
		t.Error("should pass health check after started ")
	}

	//2. test after healthCheckAged
	time.Sleep(5 * time.Second) //wait for 5 seconds
	if err := healthCheck(checkPath); err != nil {
		t.Error("should pass health check after 5 seconds", err)
	}
}

// check health by calling the http soh endpoint
func healthCheck(sohPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), healthCheckTimeout)
	defer cancel()
	msg, err := health.Check(ctx, sohPath, healthCheckTimeout)
	log.Printf("status: %s", string(msg))
	return err
}
