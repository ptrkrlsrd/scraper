package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var jobs chan ScraperTask
var logger chan string
var results map[string]ScraperResult

// ScraperTask A task containing the URL of the page you want to scrape and the delay
type ScraperTask struct {
	URL  string
	Time uint64
}

// ScraperResult The result of a scraping
type ScraperResult struct {
	ID      string
	Title   string
	Date    time.Time
	URL     string
	Content string
}

// Scrape Scrapes the given URL, and returns a ScraperResult(plus an error)
func Scrape(url string) (ScraperResult, error) {
	resp, err := http.Get(url)
	if err != nil {
		return ScraperResult{}, err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	scraperResult := ScraperResult{
		Title:   "Title",
		Date:    time.Now(),
		URL:     url,
		Content: string(bytes),
	}

	return scraperResult, nil
}

// GetScraperResult Get one ScraperResult with an ID as the param
// Example: curl localhost:4000/api/result/{key}
func GetScraperResult() func(*gin.Context) {
	return gin.HandlerFunc(func(c *gin.Context) {
		id := c.Param("id")
		scraperResult := results[id]
		c.JSON(200, scraperResult)
	})
}

// GetAllScraperResults Returns all the results
func GetAllScraperResults() func(*gin.Context) {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.JSON(200, results)
	})
}

// AddScraper Add a new scraper from JSON-data representing the ScraperTask struct
// Example data: {"URL": "https://google.com", "Time": 10}
// The URL key represents the URL you want to scrape, while the Time key represents the delay
func AddScraper() func(*gin.Context) {
	return gin.HandlerFunc(func(c *gin.Context) {
		var scraperTask ScraperTask
		c.BindJSON(&scraperTask)
		jobs <- scraperTask
		logger <- fmt.Sprintf("Added URL %s", scraperTask.URL)
	})
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	results = make(map[string]ScraperResult)
	jobs = make(chan ScraperTask)
	logger = make(chan string)

	router := gin.New()
	router.Use(gin.Logger())
	api := router.Group("/api/")
	{
		api.GET("/result/:id", GetScraperResult())
		api.GET("/results", GetAllScraperResults())
		api.POST("/scraper", AddScraper())
	}

	go func() {
		for {
			select {
			case d := <-jobs:
				go func() {
					for {
						time.Sleep(time.Duration(d.Time) * time.Second)
						scraperResult, _ := Scrape(d.URL)
						key := d.URL
						results[key] = scraperResult
						logger <- fmt.Sprintf("Scraped URL %s @ %s", scraperResult.URL, scraperResult.Date)
					}
				}()
			case logString := <-logger:
				log.Println(logString)
			}
		}
	}()

	router.Run(":4000")
}
