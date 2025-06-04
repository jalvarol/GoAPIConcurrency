package ebay

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type EbayItem struct {
	Title string `json:"title"`
	Price struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"price"`
	Image struct {
		ImageURL string `json:"imageUrl"`
	} `json:"image"`
	ItemWebURL string `json:"itemWebUrl"`
}

type EbayResponse struct {
	Items []EbayItem `json:"itemSummaries"`
}

// FetchEbayResults fetches up to 10 eBay results for the given query using the eBay Browse API.
func FetchEbayResults(query string) ([]string, error) {
	token := os.Getenv("EBAY_OAUTH_TOKEN")
	if token == "" {
		return []string{"eBay OAuth token missing"}, nil
	}
	endpoint := "https://api.ebay.com/buy/browse/v1/item_summary/search?q=" + url.QueryEscape(query) + "&limit=10"
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return []string{"eBay API request error: " + err.Error()}, nil
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []string{"eBay API error: " + err.Error()}, nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return []string{"eBay API error: status " + resp.Status}, nil
	}
	var ebayResp EbayResponse
	if err := json.NewDecoder(resp.Body).Decode(&ebayResp); err != nil {
		return []string{"eBay API decode error: " + err.Error()}, nil
	}
	var results []string
	for _, item := range ebayResp.Items {
		results = append(results, fmt.Sprintf("eBay: %s - %s %s", item.Title, item.Price.Value, item.Price.Currency))
	}
	if len(results) == 0 {
		results = append(results, "No eBay results found.")
	}
	return results, nil
}
