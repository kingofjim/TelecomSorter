package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

var ranger_army []ranger

func main() {

	ranger_army = buildArmy()
	fmt.Println(ranger_army[0])

	r := gin.Default()
	r.GET("/query/telecom", Ranger)
	r.Run(":8000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

