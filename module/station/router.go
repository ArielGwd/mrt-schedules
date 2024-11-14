package station

import (
	"net/http"

	"github.com/ArielGwd/mrt-schedules.git/module/common/response"
	"github.com/gin-gonic/gin"
)

func Initiate(router *gin.RouterGroup) {
	stationSercive := NewService()

	station := router.Group("/stations")
	station.GET("", func(ctx *gin.Context) {
		getAllStation(ctx, stationSercive)
	})
}

func getAllStation(ctx *gin.Context, srv Service) {
	data, err := srv.getAllStation()
	if err != nil {
		// handle error
		ctx.JSON(
			http.StatusBadRequest, response.ApiResponse{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			})
		return
	}

	// handle success
	ctx.JSON(
		http.StatusOK, response.ApiResponse{
			Success: true,
			Message: "Success get all station",
			Data:    data,
		})
}
