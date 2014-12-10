package main

import (
	"encoding/json"
	"fmt"
	//"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Client struct {
	APIKey string
}

func New(apikey string) *Client {
	return &Client{APIKey: apikey}
}

func (client *Client) GetDones(page int) (dones []Done, err error) {
	log.Print("GetDones =========================")

	has_next := true // First page is our next page
	var new_dones []*Done
	for next := page; has_next; next++ {
		new_dones, has_next, err = GetDonesByPage(next)
		for _, done := range new_dones {
			dones = append(dones, *done)
		}

	}
	return dones, nil
}

func GetDonesByPage(page int) (dones []*Done, has_next bool, err error) {
	log.Print("Start")
	var data struct {
		Ok       bool
		Warnings []interface{}
		Count    json.Number `json:"count,Number"`
		Next     string
		Previous string
		Results  []*Done `json:"results"`
	}

	url := fmt.Sprintf("https://idonethis.com/api/v0.1/dones/?page=%d", page)
	log.Print(url)
	http_client := &http.Client{}
	log.Printf("url: %s", url)
	// Create an http request and add a custome header
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, false, fmt.Errorf("Error formatting request: %s", err.Error())
	}

	token := os.Getenv("I_DONE_THIS_TOKEN")
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", token))

	log.Print("About to make request: %s", req)

	// Perform the request and close the connection (defered)
	res, err := http_client.Do(req)
	if err != nil {
		return nil, false, fmt.Errorf("Failed getting Dones with http err: %s", err.Error())
	}
	defer res.Body.Close()

	log.Print("About to read from response")

	// Read from body
	//body, err := ioutil.ReadAll(res.Body)

	log.Printf("Res: %s", res)
	//fmt.Printf("Body: %s", body)

	log.Print("Decoding JSON...")

	dec := json.NewDecoder(res.Body)

	log.Printf("Decoder created")
	if err := dec.Decode(&data); err != nil {
		log.Printf("ERR: %s", err.Error())
		return nil, false, fmt.Errorf("getDones failed to parse json response: %s", err.Error())
	}

	fmt.Printf("%+v\n\n", data)

	has_next = data.Next != ""

	if err != nil {
		return nil, false, err
	}

	return data.Results[:], has_next, nil // no error and make a slice of results
}

type Done struct {
	Id            json.Number            `json: "id,Number"`
	Created       string                 `json: "created"`
	Updated       string                 `json: "updated"`
	MarkedupText  string                 `json: "markedup_text"`
	Url           string                 `json: "url"`
	Team          string                 `json: "team"`
	RawText       string                 `json: "raw_text"`
	DoneDate      string                 `json: "done_date"`
	TeamShortName string                 `json:"team_short_name"`
	Owner         string                 `json:"owner"`
	Tags          []interface{}          `json: tags`
	Likes         []interface{}          `json: "likes"`
	Comments      []interface{}          `json: "comments"`
	MetaData      map[string]interface{} `json:"meta_data"`
	Permalink     string                 `json:"permalink"`
	IsGoal        bool                   `json:"is_goal"`
	GoalCompleted bool                   `json:"goal_completed"`
}

func perror(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	//url := "https://idonethis.com/api/v0.1/dones/"
	client := New("key")
	page := 1
	dones, err := client.GetDones(page)
	if err != nil {
		fmt.Errorf("Can't get dones now: %s", err)
	}

	for _, done := range dones {
		fmt.Printf("%+v\n\n", done)
	}
}
