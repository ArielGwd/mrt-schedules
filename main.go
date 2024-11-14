package main

import "github.com/gin-gonic/gin"

func main() {
	var router = gin.Default()
	router.Run(":8080")
}
