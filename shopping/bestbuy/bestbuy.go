package bestbuy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type BBProduct struct {
	Name      string  `json:"name"`
	SalePrice float64 `json:"salePrice"`
	URL       string  `json:"url"`
}

type BBResponse struct {
	Products []BBProduct `json:"products"`
}

// FetchBestBuyResults fetches shopping results from Best Buy API for a given query.
func FetchBestBuyResults(searchQuery string) ([]string, error) {
	bestBuyKey := os.Getenv("BESTBUY_API_KEY")
	if bestBuyKey == "" {
		return []string{"Best Buy API key missing"}, nil
	}
	endpoint := "https://api.bestbuy.com/v1/products((search=" + url.QueryEscape(searchQuery) + "))"
	params := "?apiKey=" + bestBuyKey + "&format=json&show=name,salePrice,url&sort=salePrice.asc&pageSize=10"
	resp, err := http.Get(endpoint + params)
	if err != nil {
		return []string{"Best Buy API error: " + err.Error()}, nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return []string{"Best Buy API error: status " + strconv.Itoa(resp.StatusCode)}, nil
	}
	var bb BBResponse
	if err := json.NewDecoder(resp.Body).Decode(&bb); err != nil {
		return []string{"Best Buy API decode error: " + err.Error()}, nil
	}
	var results []string
	for _, p := range bb.Products {
		results = append(results, fmt.Sprintf("BestBuy: %s - $%.2f", p.Name, p.SalePrice))
	}
	if len(results) == 0 {
		results = append(results, "No Best Buy results found.")
	}
	return results, nil
}
