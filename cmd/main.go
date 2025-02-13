package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Defualt()

	router.SetupRouter(r)
	fmt.Println("Server is running at 8080")
	log.Fatal(r.Run(":8080"))
}
