package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ptrkrlsrd/scraper/pkg/scraper"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	var tasks = make(chan scraper.Task) // Make an channel of ScraperTasks
	var logger = make(chan string)      // Make channel which will receive strings

	router := gin.New()             // Creates a router
	router.Use(gin.Logger())        // Add the logger middleware
	api := router.Group("/api/v1/") // Create a new API group
	{
		api.GET("/result/:id", scraper.GetResult())             // Get a result from an ID
		api.GET("/result/:id/:time", scraper.GetResultAtTime()) // Get a result from an ID and TimeStamp
		api.GET("/results", scraper.GetAllResults())            // Get all results
		api.POST("/scraper", scraper.AddScraper(tasks, logger)) // Add a new task
	}

	go scraper.Listen(tasks, logger)

	log.Println("Starting router on port 4000")
	router.Run(":4000")
}
