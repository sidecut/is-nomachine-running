package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"github.com/vmihailenco/msgpack/v5"
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
	bytes, err := msgpack.Marshal(status)
	if err != nil {
		return err
	}
	return c.Blob(http.StatusOK, "application/msgpack", bytes)
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

	// Start port 443
	go func(c *echo.Echo) {
		e.Logger.Fatal(e.StartAutoTLS(fmt.Sprintf(":%v", sslport)))
	}(e)

	// Start port 80
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", port)))
}
