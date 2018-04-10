package scraper

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// GetResult Get one ScraperResult with an ID as the param
// Example: curl localhost:4000/api/v1/result/{key}
func GetResult() func(*gin.Context) {
	return gin.HandlerFunc(func(c *gin.Context) {
		id := c.Param("id")
		scraperResults := results[id]
		c.JSON(200, scraperResults)
	})
}

// GetAllResults Returns all the results
func GetAllResults() func(*gin.Context) {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.JSON(200, results)
	})
}

// AddScraper Add a new scraper from JSON-data representing the Task struct
// Example data: {"URL": "https://google.com", "Time": 10}
// The URL key represents the URL you want to scrape, while the Time key represents the delay
func AddScraper(tasks chan Task, logger chan string) func(*gin.Context) {
	return gin.HandlerFunc(func(c *gin.Context) {
		var scraperTask Task
		c.BindJSON(&scraperTask)

		var key = md5Hash(scraperTask.URL)
		scraperTask.Key = key

		tasks <- scraperTask
		logger <- fmt.Sprintf("Added URL %s", scraperTask.URL)
		c.String(200, scraperTask.Key)
	})
}
