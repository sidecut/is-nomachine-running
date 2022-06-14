package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/requestid"
	"github.com/gin-contrib/static"
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

	e.Use(gzip.Gzip(gzip.DefaultCompression))
	// e.Use(middleware.Logger())
	e.Use(requestid.New())

	e.Use(cors.Default())

	e.GET("/api", statusAPI)
	e.Use(static.Serve("/", static.LocalFile("dist", true)))

	viper.AutomaticEnv()
	viper.SetEnvPrefix("isno")
	port := viper.GetInt("port")

	// Start port 80
	e.Run(fmt.Sprintf(":%v", port))
}
