package scraper

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/ptrkrlsrd/scraper/pkg/title"
)

type Service struct {
	Results map[string][]Result
}

func NewService() Service {
	results := make(map[string][]Result)
	return Service{
		Results: results,
	}
}

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

// init Init
func init() {
	log.SetPrefix("[Scraper] ")
}

// md5Hash Run MD5-hashing on a string
func md5Hash(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Scrape Scrapes the given URL, and returns a Result(plus an error)
func (scraperTask *Task) Scrape() (Result, error) {
	resp, err := http.Get(scraperTask.URL)
	if err != nil {
		return Result{}, err
	}

	body := resp.Body
	key := md5Hash(scraperTask.URL)
	bytes, err := ioutil.ReadAll(body)
	pageTitle, _ := title.GetHtmlTitle(body)

	if err != nil {
		return Result{}, err
	}

	scraperResult := Result{
		ID:      key,
		Title:   pageTitle,
		Date:    time.Now(),
		URL:     scraperTask.URL,
		Content: string(bytes),
	}

	return scraperResult, nil
}

// Listen Listen takes a chan of Tasks and a chan of strings and listens for events
func (service *Service) Listen(tasks chan Task, logger chan string) {
	for {
		select {
		case task := <-tasks:
			go func() {
				for {
					time.Sleep(time.Duration(task.Time) * time.Second)
					scraperResult, err := task.Scrape()
					if err != nil {
						log.Println(err)
					}
					// Probably not safe for concurrency
					service.Results[task.Key] = append(service.Results[task.Key], scraperResult)
					logger <- fmt.Sprintf("Scraped URL %s @ %s", scraperResult.URL, scraperResult.Date)
				}
			}()
		case logString := <-logger:
			log.Println(logString)
		}
	}
}
