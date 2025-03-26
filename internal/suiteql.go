package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/TylerBrock/colorjson"
	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
)

type Credentials struct {
	AccountID      string
	ConsumerKey    string
	ConsumerSecret string
	Token          string
	TokenSecret    string
	Realm          string
}

type SuiteQLClient struct {
	creds *Credentials
}

func NewSuiteQLClient() (*SuiteQLClient, error) {
	creds, err := getCredentials()
	if err != nil {
		return nil, err
	}
	return &SuiteQLClient{creds: creds}, nil
}

func (c *SuiteQLClient) ExecuteQuery(query string, limit, offset *int) (string, error) {
	config := oauth1.NewConfig(c.creds.ConsumerKey, c.creds.ConsumerSecret)
	config.Realm = c.creds.Realm
	config.Signer = &oauth1.HMAC256Signer{ConsumerSecret: c.creds.ConsumerSecret}

	token := oauth1.NewToken(c.creds.Token, c.creds.TokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	baseURL := fmt.Sprintf("https://%s.suitetalk.api.netsuite.com/services/rest/query/v1/suiteql", c.creds.AccountID)

	// Parse the base URL
	urlBuilder, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("error parsing URL: %v", err)
	}

	// Create query parameters
	queryParams := url.Values{}
	if limit != nil {
		queryParams.Set("limit", fmt.Sprintf("%d", *limit))
	}
	if offset != nil {
		queryParams.Set("offset", fmt.Sprintf("%d", *offset))
	}

	// Set the query parameters
	urlBuilder.RawQuery = queryParams.Encode()

	// Prepare request body
	payload := map[string]string{
		"q": query,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Create request
	req, err := http.NewRequest("POST", urlBuilder.String(), strings.NewReader(string(jsonData)))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("prefer", "transient")

	// Send request
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	// Check for non-200 status codes
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("error: received status code %d: %s", resp.StatusCode, string(body))
	}

	// Parse JSON for pretty printing with colors
	var obj any
	if err := json.Unmarshal(body, &obj); err != nil {
		return "", fmt.Errorf("error parsing JSON response: %v", err)
	}

	// Create a color formatter
	formatter := colorjson.NewFormatter()
	formatter.Indent = 2

	// Format the JSON with colors
	output, err := formatter.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("error formatting JSON: %v", err)
	}

	return string(output), nil
}

func getCredentials() (*Credentials, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	creds := &Credentials{
		AccountID:      os.Getenv("NETSUITE_ACCOUNT_ID"),
		ConsumerKey:    os.Getenv("NETSUITE_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("NETSUITE_CONSUMER_SECRET"),
		Token:          os.Getenv("NETSUITE_TOKEN"),
		TokenSecret:    os.Getenv("NETSUITE_TOKEN_SECRET"),
		Realm:          os.Getenv("NETSUITE_ACCOUNT_ID"),
	}

	// Validate all required credentials are present
	missing := []string{}
	if creds.AccountID == "" {
		missing = append(missing, "NETSUITE_ACCOUNT_ID")
	}
	if creds.ConsumerKey == "" {
		missing = append(missing, "NETSUITE_CONSUMER_KEY")
	}
	if creds.ConsumerSecret == "" {
		missing = append(missing, "NETSUITE_CONSUMER_SECRET")
	}
	if creds.Token == "" {
		missing = append(missing, "NETSUITE_TOKEN")
	}
	if creds.TokenSecret == "" {
		missing = append(missing, "NETSUITE_TOKEN_SECRET")
	}
	if creds.Realm == "" {
		missing = append(missing, "NETSUITE_ACCOUNT_ID")
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s",
			strings.Join(missing, ", "))
	}

	return creds, nil
}
