package main

import (
	"github.com/ArielGwd/mrt-schedules.git/module/station"
	"github.com/gin-gonic/gin"
)

func main() {
	InitiateRoutes()
}

func InitiateRoutes() {
	var (
		router = gin.Default()
		api    = router.Group("/v1/api/")
	)

	station.Initiate(api)
	router.Run(":8080")
	print("Server is running on port localhost:8080")
}
