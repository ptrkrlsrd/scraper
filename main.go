package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ptrkrlsrd/scraper/pkg/title"
)

var results map[string]map[time.Time]ScraperResult

// ScraperTask A task containing the URL of the page you want to scrape and the delay
type ScraperTask struct {
	Key  string
	URL  string `json:"url"`
	Time uint64 `json:"time"`
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

	key := md5Hash(url)
	bytes, err := ioutil.ReadAll(resp.Body)
	pageTitle, _ := title.GetHtmlTitle(resp.Body)

	if err != nil {
		return ScraperResult{}, err
	}

	scraperResult := ScraperResult{
		ID:      key,
		Title:   pageTitle,
		Date:    time.Now(),
		URL:     url,
		Content: string(bytes),
	}

	return scraperResult, nil
}

// GetScraperResult Get one ScraperResult with an ID as the param
// Example: curl localhost:4000/api/v1/result/{key}
func GetScraperResult() func(*gin.Context) {
	return gin.HandlerFunc(func(c *gin.Context) {
		id := c.Param("id")
		scraperResults := results[id]
		c.JSON(200, scraperResults)
	})
}

// GetScraperResultAtTime Get one ScraperResult with an ID as the param
// Example: curl localhost:4000/api/v1/result/{key}
func GetScraperResultAtTime() func(*gin.Context) {
	return gin.HandlerFunc(func(c *gin.Context) {
		id := c.Param("id")
		timeStamp := c.Param("time")
		t, _ := time.Parse(time.RFC3339, timeStamp)
		scraperResult := results[id][t]
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
func AddScraper(tasks chan ScraperTask, logger chan string) func(*gin.Context) {
	return gin.HandlerFunc(func(c *gin.Context) {
		var scraperTask ScraperTask
		c.BindJSON(&scraperTask)

		var key = md5Hash(scraperTask.URL)
		scraperTask.Key = key

		tasks <- scraperTask
		logger <- fmt.Sprintf("Added URL %s", scraperTask.URL)
		c.String(200, scraperTask.Key)
	})
}

func md5Hash(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	var tasks = make(chan ScraperTask) // Make an channel of ScraperTasks
	var logger = make(chan string)     // Make channel which will receive strings
	results = make(map[string]map[time.Time]ScraperResult)

	router := gin.New()             // Creates a router
	router.Use(gin.Logger())        // Add the logger middleware
	api := router.Group("/api/v1/") // Create a new API group
	{
		api.GET("/result/:id", GetScraperResult())             // Get a result from an ID
		api.GET("/result/:id/:time", GetScraperResultAtTime()) // Get a result from an ID and TimeStamp
		api.GET("/results", GetAllScraperResults())            // Get all results
		api.POST("/scraper", AddScraper(tasks, logger))        // Add a new task
	}

	go func() {
		for {
			select {
			case d := <-tasks:
				go func() {
					for {
						time.Sleep(time.Duration(d.Time) * time.Second)
						scraperResult, _ := Scrape(d.URL)
						results[d.Key] = map[time.Time]ScraperResult{time.Now(): scraperResult}
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
