package main

import (
	"github.com/gin-gonic/gin"
)

var ranger_army []ranger

func main() {

	ranger_army = buildArmy()

	r := gin.Default()
	r.GET("/query/telecom", Ranger)
	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

