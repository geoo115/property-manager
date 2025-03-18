package main

import (
	"fmt"
	"log"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/events"
	"github.com/geoo115/property-manager/router"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database and Redis
	db.Init()
	db.InitRedis()

	// Initialize Kafka (Producer)
	events.InitKafka()

	// ✅ Start Kafka Consumer in a separate goroutine
	fmt.Println("🚀 Starting Kafka Consumer...")
	go events.StartKafkaConsumer()

	// Setup Gin router
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup API routes
	router.SetupRouter(r)

	fmt.Println("🚀 Server is running on port 8080")
	log.Fatal(r.Run(":8080"))
}
