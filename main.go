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

type ScraperTask struct {
	URL  string
	Time uint64
}

type ScraperResult struct {
	ID      string
	Title   string
	Date    time.Time
	URL     string
	Content string
}

func Scrape(url string) (ScraperResult, error) {
	resp, err := http.Get(url)
	if err != nil {
		return ScraperResult{}, err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	scraperResult := ScraperResult{Title: "Title", Date: time.Now(), URL: url, Content: string(bytes)}

	return scraperResult, nil
}

func GetScraperResult() func(*gin.Context) {
	return gin.HandlerFunc(func(c *gin.Context) {
		id := c.Param("id")
		scraperResult := results[id]
		c.JSON(200, scraperResult)
	})
}

func GetAllScraperResults() func(*gin.Context) {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.JSON(200, results)
	})
}

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
						key := "key" //scraperResult.URL
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
