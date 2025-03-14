package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sqc157400661/helper/conf"
	"net/http"
	"time"
)

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:              conf.GetStringSlice("corsDomains"),
		AllowMethods:              []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:              []string{"Origin", "Authorization", "Content-Type", "Content-Length", "Access-Token", "Buservice-Id", "access-token", "Accept", "Engine"},
		ExposeHeaders:             []string{"Content-Type", "Content-Length", "Content-Language", "Access-Control-Allow-Origin", "Cache-Control", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods"},
		AllowWildcard:             true,
		AllowPrivateNetwork:       true,
		AllowCredentials:          true,
		OptionsResponseStatusCode: http.StatusOK,
		MaxAge:                    12 * time.Hour,
	})
}
