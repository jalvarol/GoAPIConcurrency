package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/freshman-tech/news-demo-starter-files/news"

	"github.com/joho/godotenv"
)

var tpl = template.Must(template.ParseFiles("index.html"))

type Search struct {
	Query           string
	NextPage        int
	TotalPages      int
	Results         *news.Results
	ShoppingResults []string // Placeholder for shopping results
}

func (s *Search) IsLastPage() bool {
	return s.NextPage >= s.TotalPages
}
func (s *Search) CurrentPage() int {
	if s.NextPage == 1 {
		return s.NextPage
	}

	return s.NextPage - 1
}
func (s *Search) PreviousPage() int {
	return s.CurrentPage() - 1
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}
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

		// Fetch news concurrently
		go func() {
			results, err := newsapi.FetchEverything(searchQuery, page)
			if err != nil {
				newsErrCh <- err
				return
			}
			newsCh <- results
		}()

		// Fetch shopping results concurrently (Best Buy API)
		go func() {
			bestBuyKey := os.Getenv("BESTBUY_API_KEY")
			if bestBuyKey == "" {
				shoppingCh <- []string{"Best Buy API key missing"}
				return
			}
			// Example: Search for products using Best Buy API
			// Docs: https://developer.bestbuy.com/documentation/products-api
			endpoint := "https://api.bestbuy.com/v1/products((search=" + url.QueryEscape(searchQuery) + "))"
			params := "?apiKey=" + bestBuyKey + "&format=json&show=name,salePrice,url&sort=salePrice.asc&pageSize=3"
			resp, err := http.Get(endpoint + params)
			if err != nil {
				shoppingCh <- []string{"Best Buy API error: " + err.Error()}
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				shoppingCh <- []string{"Best Buy API error: status " + strconv.Itoa(resp.StatusCode)}
				return
			}
			type BBProduct struct {
				Name      string  `json:"name"`
				SalePrice float64 `json:"salePrice"`
				URL       string  `json:"url"`
			}
			type BBResponse struct {
				Products []BBProduct `json:"products"`
			}
			var bb BBResponse
			if err := json.NewDecoder(resp.Body).Decode(&bb); err != nil {
				shoppingCh <- []string{"Best Buy API decode error: " + err.Error()}
				return
			}
			var results []string
			for _, p := range bb.Products {
				results = append(results, fmt.Sprintf("BestBuy: %s - $%.2f", p.Name, p.SalePrice))
			}
			if len(results) == 0 {
				results = append(results, "No Best Buy results found.")
			}
			shoppingCh <- results
		}()

		var results *news.Results
		var shoppingResults []string
		var fetchErr error
		for i := 0; i < 2; i++ {
			select {
			case r := <-newsCh:
				results = r
			case err := <-newsErrCh:
				fetchErr = err
			case s := <-shoppingCh:
				shoppingResults = s
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

		search := &Search{
			Query:           searchQuery,
			NextPage:        nextPage,
			TotalPages:      int(math.Ceil(float64(results.TotalResults) / float64(newsapi.PageSize))),
			Results:         results,
			ShoppingResults: shoppingResults,
		}
		if ok := !search.IsLastPage(); ok {
			search.NextPage++
		}
		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf.WriteTo(w)
		fmt.Printf("%+v", results)
		fmt.Printf("%+v", results)
		fmt.Println("Search Query is: ", searchQuery)
		fmt.Println("Page is: ", page)
	}
}
func main() {
	err := godotenv.Load()
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

	fs := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("/search", searchHandler(newsapi))
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux)
}
