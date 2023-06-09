package main

import (
	"github.com/bmbbms/alert-rlist/pkg/influxdbAlert"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	engine.GET("/alert", influxdbAlert.GetAlertInfo)
	engine.Run(":8080")
}
