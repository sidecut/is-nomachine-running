package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("port", 80)
	viper.SetDefault("sslport", 443)
}

func statusAPI(c *gin.Context) {
	status, err := getStatus()
	// err = errors.New("Something bad happened")
	if err != nil {
		// TODO: log this error
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, status)
}

func main() {
	e := gin.Default()

	// e := echo.New()
	// e.Use(middleware.Gzip())
	// e.Use(middleware.Logger())
	// e.Use(middleware.RequestID())

	// corsConfig := middleware.CORSConfig{AllowOrigins: []string{"*"}}
	// e.Use(middleware.CORSWithConfig(corsConfig))

	e.GET("/api", statusAPI)
	e.NoRoute(func(ctx *gin.Context) {
		e.StaticFS("/", gin.Dir("dist", true))
	})

	viper.AutomaticEnv()
	viper.SetEnvPrefix("isno")
	port := viper.GetInt("port")
	// sslport := viper.GetInt("sslport")

	// // Start port 443
	// go func(c *echo.Echo) {
	// 	e.Logger.Fatal(e.StartAutoTLS(fmt.Sprintf(":%v", sslport)))
	// }(e)

	// Start port 80
	e.Run(fmt.Sprintf(":%v", port))
}
