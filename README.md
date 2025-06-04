# GoHeadlines

This project demonstrates Go concurrency by fetching news articles and shopping partner results in parallel, then displaying them side-by-side in a modern web UI.

## Features
- Search for news articles by keyword (powered by NewsAPI)
- Compare shopping results from multiple partners (Best Buy, more coming soon)
- Modern, responsive UI inspired by the San Jose Sharks color palette
- Built with Go, HTML, and CSS
- Clear use of goroutines and channels for concurrent API calls

## Prerequisites
- [Go](https://golang.org/dl/) (tested with 1.15+)
- NewsAPI key ([get one here](https://newsapi.org/))
- (Optional) Best Buy API key ([get one here](https://developer.bestbuy.com/))

## Setup
1. Clone this repository:
   ```sh
   git clone https://github.com/jalvarol/GoAPIConcurrency.git
   cd GoAPIConcurrency/GoAPIConcurrency
   ```
2. Create a `.env` file in the project root:
   ```env
   PORT=3000
   NEWS_API_KEY=your_news_api_key_here
   BESTBUY_API_KEY=your_bestbuy_api_key_here
   ```
   Replace the values with your actual API keys.

3. Run the app:
   ```sh
   go run main.go
   ```
   Then open [http://localhost:3000](http://localhost:3000) in your browser.

## How it works
- The backend uses goroutines and channels to fetch news and shopping results concurrently.
- The UI displays news on the left and shopping results on the right.
- Easily extensible to add more shopping APIs (Amazon, Target, Walmart, etc).

## License
MIT
