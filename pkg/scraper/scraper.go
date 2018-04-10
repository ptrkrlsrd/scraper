package scraper

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ptrkrlsrd/scraper/pkg/title"
)

var results map[string]map[time.Time]Result

// Task A task containing the URL of the page you want to scrape and the delay
type Task struct {
	Key  string
	URL  string `json:"url"`
	Time uint64 `json:"time"`
}

// Result The result of a scraping
type Result struct {
	ID      string
	Title   string
	Date    time.Time
	URL     string
	Content string
}

// Scrape Scrapes the given URL, and returns a Result(plus an error)
func Scrape(url string) (Result, error) {
	resp, err := http.Get(url)
	if err != nil {
		return Result{}, err
	}

	key := md5Hash(url)

	pageTitle, _ := title.GetHtmlTitle(resp.Body)
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Result{}, err
	}

	scraperResult := Result{
		ID:      key,
		Title:   pageTitle,
		Date:    time.Now(),
		URL:     url,
		Content: string(bytes),
	}

	return scraperResult, nil
}

// GetResult Get one ScraperResult with an ID as the param
// Example: curl localhost:4000/api/v1/result/{key}
func GetResult() func(*gin.Context) {
	return gin.HandlerFunc(func(c *gin.Context) {
		id := c.Param("id")
		scraperResults := results[id]
		c.JSON(200, scraperResults)
	})
}

// GetResultAtTime Get one ScraperResult with an ID as the param
// Example: curl localhost:4000/api/v1/result/{key}
func GetResultAtTime() func(*gin.Context) {
	return gin.HandlerFunc(func(c *gin.Context) {
		id := c.Param("id")
		timeStamp := c.Param("time")
		t, _ := time.Parse(time.RFC3339, timeStamp)
		scraperResult := results[id][t]
		c.JSON(200, scraperResult)
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

func md5Hash(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}
