package news

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	HTTPClient *http.Client
	APIKey     string
	PageSize   int
}

type Results struct {
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

type Article struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
	Source      Source    `json:"source"`
}

type Source struct {
	Name string `json:"name"`
}

func (a *Article) FormatPublishedDate() string {
	return a.PublishedAt.Format("Jan 2, 2006 15:04")
}

func NewClient(httpClient *http.Client, apiKey string, pageSize int) *Client {
	return &Client{
		HTTPClient: httpClient,
		APIKey:     apiKey,
		PageSize:   pageSize,
	}
}

func (c *Client) FetchEverything(query, page string) (*Results, error) {
	endpoint := "https://newsapi.org/v2/everything"
	params := url.Values{}
	params.Add("q", query)
	params.Add("pageSize", strconv.Itoa(c.PageSize))
	params.Add("page", page)
	params.Add("apiKey", c.APIKey)

	url := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("newsapi: unexpected status code: %d", resp.StatusCode)
	}

	var results Results
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}
	return &results, nil
}
