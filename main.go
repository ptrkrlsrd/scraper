package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ptrkrlsrd/scraper/pkg/scraper"
)

var results map[string]map[time.Time]scraper.Result

func main() {
	gin.SetMode(gin.ReleaseMode)

	var tasks = make(chan scraper.Task) // Make an channel of ScraperTasks
	var logger = make(chan string)      // Make channel which will receive strings
	results = make(map[string]map[time.Time]scraper.Result)

	router := gin.New()             // Creates a router
	router.Use(gin.Logger())        // Add the logger middleware
	api := router.Group("/api/v1/") // Create a new API group
	{
		api.GET("/result/:id", scraper.GetResult())             // Get a result from an ID
		api.GET("/result/:id/:time", scraper.GetResultAtTime()) // Get a result from an ID and TimeStamp
		api.GET("/results", scraper.GetAllResults())            // Get all results
		api.POST("/scraper", scraper.AddScraper(tasks, logger)) // Add a new task
	}

	go func() {
		for {
			select {
			case d := <-tasks:
				go func() {
					for {
						time.Sleep(time.Duration(d.Time) * time.Second)
						scraperResult, _ := scraper.Scrape(d.URL)
						results[d.Key] = map[time.Time]scraper.Result{time.Now(): scraperResult}
						logger <- fmt.Sprintf("Scraped URL %s @ %s", scraperResult.URL, scraperResult.Date)
					}
				}()
			case logString := <-logger:
				log.Println(logString)
			}
		}
	}()

	log.Println("Starting router on port 4000")
	router.Run(":4000")
}
