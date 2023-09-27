package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
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

	e.Logger.SetLevel(log.INFO)
	e.Logger.Infof("*** STARTING PID %v", os.Getpid())

	// Start port 443
	go func(c *echo.Echo) {
		e.Logger.Fatal(e.StartAutoTLS(fmt.Sprintf(":%v", sslport)))
	}(e)

	// Start port 80
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", port)))
}
