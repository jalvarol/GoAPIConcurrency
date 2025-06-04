// main.go - Entry point for GoHeadlines web application
// Handles HTTP server setup, routing, and concurrent fetching of news and shopping results
package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/jalvarol/goheadlines/news"
	"github.com/jalvarol/goheadlines/shopping/bestbuy"
	"github.com/jalvarol/goheadlines/shopping/ebay"
	"github.com/joho/godotenv"
)

// tpl is the parsed HTML template for rendering the main page
var tpl = template.Must(template.ParseFiles("index.html"))

// Search holds data for a search query, including pagination and results
// Used to pass data to the HTML template
// Query: the search keyword
// NextPage: the next page number for pagination
// TotalPages: total number of pages available
// Results: news API results
// ShoppingResults: shopping API results (Best Buy, etc)
type Search struct {
	Query           string
	NextPage        int
	TotalPages      int
	Results         *news.Results
	ShoppingResults []string // Placeholder for shopping results
	EbayResults     []string // Placeholder for eBay results
}

// IsLastPage returns true if the current page is the last page of results
func (s *Search) IsLastPage() bool {
	return s.NextPage >= s.TotalPages
}

// CurrentPage returns the current page number
func (s *Search) CurrentPage() int {
	if s.NextPage == 1 {
		return s.NextPage
	}

	return s.NextPage - 1
}

// PreviousPage returns the previous page number
func (s *Search) PreviousPage() int {
	return s.CurrentPage() - 1
}

// indexHandler renders the main search page (no results)
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

// searchHandler handles /search requests, fetches news and shopping results concurrently, and renders the results page
func searchHandler(newsapi *news.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		params := u.Query()
		searchQuery := params.Get("q")
		page := params.Get("page")
		if page == "" {
			page = "1"
		}

		// Channels for concurrency
		newsCh := make(chan *news.Results)
		newsErrCh := make(chan error)
		shoppingCh := make(chan []string)
		ebayCh := make(chan []string)
		newsTimerCh := make(chan time.Duration)
		shoppingTimerCh := make(chan time.Duration)
		ebayTimerCh := make(chan time.Duration)

		// Fetch news concurrently
		go func() {
			start := time.Now()
			results, err := newsapi.FetchEverything(searchQuery, page)
			newsTimerCh <- time.Since(start)
			if err != nil {
				newsErrCh <- err
				return
			}
			newsCh <- results
		}()

		// Fetch shopping results concurrently (Best Buy API)
		go func() {
			start := time.Now()
			results, err := bestbuy.FetchBestBuyResults(searchQuery)
			shoppingTimerCh <- time.Since(start)
			if err != nil {
				shoppingCh <- []string{"Best Buy API error: " + err.Error()}
				return
			}
			shoppingCh <- results
		}()

		// Fetch eBay results concurrently
		go func() {
			start := time.Now()
			results, err := ebay.FetchEbayResults(searchQuery)
			ebayTimerCh <- time.Since(start)
			if err != nil {
				ebayCh <- []string{"eBay API error: " + err.Error()}
				return
			}
			ebayCh <- results
		}()

		// Wait for all goroutines to finish and collect results/errors
		var results *news.Results
		var shoppingResults, ebayResults []string
		var fetchErr error
		var newsElapsed, shoppingElapsed, ebayElapsed time.Duration
		for i := 0; i < 6; i++ {
			select {
			case r := <-newsCh:
				results = r
			case err := <-newsErrCh:
				fetchErr = err
			case s := <-shoppingCh:
				shoppingResults = s
			case e := <-ebayCh:
				ebayResults = e
			case t := <-newsTimerCh:
				newsElapsed = t
			case t := <-shoppingTimerCh:
				shoppingElapsed = t
			case t := <-ebayTimerCh:
				ebayElapsed = t
			}
		}
		if fetchErr != nil {
			http.Error(w, fetchErr.Error(), http.StatusInternalServerError)
			return
		}

		nextPage, err := strconv.Atoi(page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Prepare data for template rendering
		search := &Search{
			Query:           searchQuery,
			NextPage:        nextPage,
			TotalPages:      int(math.Ceil(float64(results.TotalResults) / float64(newsapi.PageSize))),
			Results:         results,
			ShoppingResults: shoppingResults,
			EbayResults:     ebayResults,
		}
		if ok := !search.IsLastPage(); ok {
			search.NextPage++
		}
		// Format durations as ms strings
		newsMs := fmt.Sprintf("%.2f ms", float64(newsElapsed.Microseconds())/1000)
		shoppingMs := fmt.Sprintf("%.2f ms", float64(shoppingElapsed.Microseconds())/1000)
		ebayMs := fmt.Sprintf("%.2f ms", float64(ebayElapsed.Microseconds())/1000)
		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, struct {
			*Search
			NewsElapsed     string
			ShoppingElapsed string
			EbayElapsed     string
		}{
			Search:          search,
			NewsElapsed:     newsMs,
			ShoppingElapsed: shoppingMs,
			EbayElapsed:     ebayMs,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf.WriteTo(w)
		fmt.Printf("%+v", results)
		fmt.Println("Search Query is: ", searchQuery)
		fmt.Println("Page is: ", page)
		fmt.Printf("News API time: %v\n", newsElapsed)
		fmt.Printf("Shopping API time: %v\n", shoppingElapsed)
		fmt.Printf("eBay API time: %v\n", ebayElapsed)
	}
}

// main initializes environment, sets up the HTTP server, and starts listening for requests
func main() {
	err := godotenv.Load() // Load environment variables from .env file (if present)
	if err != nil {
		log.Println("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		log.Fatal("Env: apiKey must be set")
	}

	myClient := &http.Client{Timeout: 10 * time.Second}
	newsapi := news.NewClient(myClient, apiKey, 20)

	mux := http.NewServeMux()

	// Serve static assets (CSS, images, etc)
	fs := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// Route handlers
	mux.HandleFunc("/search", searchHandler(newsapi))
	mux.HandleFunc("/", indexHandler)

	http.ListenAndServe(":"+port, mux)
}
