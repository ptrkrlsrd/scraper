# Scraper
## A scheduled scraper written in Go


## Usage
* Add a new task
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


## Lines of code
```
-------------------------------------------------------------------------------
Language                     files          blank        comment           code
-------------------------------------------------------------------------------
Go                               1             15              9             90
-------------------------------------------------------------------------------
```
