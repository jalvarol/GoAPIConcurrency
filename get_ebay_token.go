// get_ebay_token.go
// Run this file with: go run get_ebay_token.go
// It will print your eBay OAuth2 application access token for use in your .env file.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	start := time.Now()
	appID := os.Getenv("EBAY_APP_ID")
	certID := os.Getenv("EBAY_CERT_ID")
	if appID == "" || certID == "" {
		fmt.Println("EBAY_APP_ID and EBAY_CERT_ID must be set in your .env file.")
		return
	}

	endpoint := "https://api.ebay.com/identity/v1/oauth2/token"
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "https://api.ebay.com/oauth/api_scope")

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(appID, certID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("HTTP error:", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		fmt.Println("eBay token error:", string(body))
		return
	}
	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Decode error:", err)
		return
	}
	fmt.Println("Your eBay OAuth2 token (add to .env as EBAY_OAUTH_TOKEN):")
	fmt.Println(result.AccessToken)
	fmt.Printf("(Expires in %d seconds)\n", result.ExpiresIn)
	fmt.Printf("Token fetch took: %v\n", time.Since(start))
}
