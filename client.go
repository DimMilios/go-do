package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

/*
	List
		GET /1/lists/[idList]/cards - Get an array of Cards on a List
		POST /1/cards - Create a new Card on a List

	Card
		PUT /1/cards/[card id or shortlink] - Update the contents of a Card
		POST /1/cards/[card id or shortlink]/actions/comments - Add a comment to a Card
*/

type Board struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	ShortURL string `json:"shortUrl"`
}

type Card struct {
	ID        string   `json:"id"`
	Desc      string   `json:"desc"`
	Due       string   `json:"due"`
	Email     string   `json:"email"`
	IDBoard   string   `json:"idBoard"`
	IDList    string   `json:"idList"`
	IDMembers []string `json:"idMembers"`
	Name      string   `json:"name"`
	ShortURL  string   `json:"shortUrl"`
	URL       string   `json:"url"`
}

type Client struct {
	apiKey     string
	token      string
	baseURL    string
	HTTPClient *http.Client
}

func NewClient(apiKey, token string) *Client {
	return &Client{
		apiKey: apiKey,
		token:  token,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
		baseURL: "https://api.trello.com/1",
	}
}

func (c *Client) createURL(path string, queryParams string) *url.URL {
	u := &url.URL{
		Scheme:   "https",
		Host:     "api.trello.com",
		Path:     "1/" + path,
		RawQuery: fmt.Sprintf("%s&key=%s&token=%s", queryParams, c.apiKey, c.token),
	}
	return u
}

func (c *Client) demoCall(qURL *url.URL) {
	// trelloURL := fmt.Sprintf("%s/members/me/boards?fields=name,url&key=%s&token=%s", c.baseURL, c.apiKey, c.token)
	req, _ := http.NewRequest("GET", qURL.String(), nil)
	req.Header.Set("Content-Type", "application/json")
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}

	var result []interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("%s\n", body)
		log.Fatal(err)
	}
	fmt.Printf("Response body: %s\n", result)
}

// POST https://api.trello.com/1/boards/?name={name}&key=APIKey&token=APIToken
func (c *Client) createBoardHandler(name string) error {
	// trelloURL := fmt.Sprintf("%s/members/me/boards?fields=name,url&key=%s&token=%s", c.baseURL, c.apiKey, c.token)
	qURL := c.createURL("boards/", fmt.Sprintf("name=%s", url.QueryEscape(name)))

	req, _ := http.NewRequest("POST", qURL.String(), nil)
	req.Header.Set("Content-Type", "application/json")
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Print(err)
		return err
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		log.Printf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Print(err)
		return err
	}

	var result []interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("%s\n", body)
		log.Print(err)
		return err
	}
	fmt.Printf("Response body: %s\n", result)
	fmt.Printf("Board %s successfully created!\n", name)

	return nil
}
