package twitter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/charmbracelet/huh"
)

const (
	postTweetURL = "https://api.twitter.com/2/tweets"
)

func (c *Client) PostTweet() error {
	var tweet string
	huh.NewText().
		Title("Write your tweet").
		Value(&tweet).
		CharLimit(240).
		Validate(func(str string) error {
			if str == "" {
				return fmt.Errorf("tweet can't be empty")
			}
			return nil
		}).
		Run()

	err := c.sendTweet(tweet)

	return err
}

func (c *Client) sendTweet(tweet string) error {
	fmt.Printf("Tweet sent: %s\n", tweet)

	// Construct tweet data
	tweetData := struct {
		Text string `json:"text"`
	}{
		Text: tweet,
	}

	// Encode tweet data to JSON
	jsonData, err := json.Marshal(tweetData)
	if err != nil {
		return fmt.Errorf("error encoding tweet data: %v", err)
	}

	//JSON Body
	body := []byte(jsonData)
	// create a new request
	r, errn := http.NewRequest("POST", postTweetURL, bytes.NewBuffer(body))

	if errn != nil {
		return errn
	}

	//add headers
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.accessToken.AccessToken))

	// create a client
	client := &http.Client{}

	res, errd := client.Do(r)

	if errd != nil {
		return errd
	}

	defer res.Body.Close()

	// handling the response

	tweetResponse := &TweetResponse{}

	derr := json.NewDecoder(res.Body).Decode(tweetResponse)

	if derr != nil {
		return derr
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf(res.Status)
	}

	fmt.Print(tweetResponse)

	return nil
}
