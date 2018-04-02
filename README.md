# Scraper
## A scheduled scraper written in Go


## Usage
* Add a new task
    ```
    curl -X POST localhost:4000/api/scraper -d '{"URL": "https://theverge.com", "Time": 10}'
    ```
* Get all results
    ```
    curl localhost:4000/api/results
    ```
