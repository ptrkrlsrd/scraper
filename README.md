# Scraper
## A scheduled scraper written in Go


## Usage
* Add a new task
    * Note: The time is in seconds
    ```
    curl -X POST localhost:4000/api/v1/scraper -d '{"URL": "https://theverge.com", "Time": 10}'
    ```
* Get all results
    ```
    curl localhost:4000/api/v1/results
    ```

* Get a result
    ```
    curl localhost:4000/api/v1/result/{id}
    ```
