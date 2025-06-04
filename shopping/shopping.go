package shopping

import (
	"github.com/jalvarol/goheadlines/shopping/bestbuy"
)

// FetchBestBuyResults fetches shopping results from Best Buy API for a given query.
// Deprecated: Use shopping/bestbuy package instead.
func FetchBestBuyResults(searchQuery string) ([]string, error) {
	return bestbuy.FetchBestBuyResults(searchQuery)
}
