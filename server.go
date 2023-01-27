package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	status "github.com/sidecut/is-nomachine-running/gen/proto/go/nomachine_status/v1"
)

func init() {
	viper.SetDefault("port", 80)
	viper.SetDefault("sslport", 443)
}

func statusAPI(c echo.Context) error {
	status, err := getStatus()
	if err != nil {
		// TODO: log this error
		return err
	}
	return c.JSON(http.StatusOK, status)
}

func main() {
	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())

	corsConfig := middleware.CORSConfig{AllowOrigins: []string{"*"}}
	e.Use(middleware.CORSWithConfig(corsConfig))

	e.GET("/api", statusAPI)
	e.Static("/", "dist")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("isno")
	port := viper.GetInt("port")
	sslport := viper.GetInt("sslport")

	go func(c *echo.Echo) {
		listenOn := "127.0.0.1:8080"
		listener, err := net.Listen("tcp", listenOn)
		if err != nil {
			e.Logger.Fatal(fmt.Errorf("failed to listen on %s: %w", listenOn, err))
		}

		server := grpc.NewServer()
		status.RegisterGetStatusResponseServiceServer(server, &statusServiceServer{})
		log.Println("Listening on", listenOn)
		if err := server.Serve(listener); err != nil {
			e.Logger.Fatal(fmt.Errorf("failed to serve gRPC server: %w", err))
		}
	}(e)

	// Start port 443
	go func(c *echo.Echo) {
		e.Logger.Fatal(e.StartAutoTLS(fmt.Sprintf(":%v", sslport)))
	}(e)

	// Start port 80
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", port)))
}

type statusServiceServer struct {
	status.UnimplementedGetStatusResponseServiceServer
}

func (s *statusServiceServer) GetStatus(ctx context.Context, req *status.GetStatusRequest) (*status.GetStatusResponse, error) {
	return &status.GetStatusResponse{
		HostName:         "TEST-HOST-NANE",
		NoMachineRunning: true,
		ClientAttached:   true,
	}, nil
}
