package twitter

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

const (
	generateAccessTokenUrl = "https://api.twitter.com/2/oauth2/token"
)

type TweetResponse struct {
	Data struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"data"`
}

type AccessTokenResponse struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
}

var logger = log.NewWithOptions(os.Stderr, log.Options{
	Prefix: "Cooking 🍪 ",
})

func (c *Client) AuthorizeAccount() error {
	logger.Info("Opening Browser for account authorization")
	// Define the parameters
	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", c.clientId)
	params.Set("redirect_uri", "http://localhost:3003/token")
	params.Set("scope", "tweet.read tweet.write users.read follows.read follows.write")
	params.Set("state", "state")
	params.Set("code_challenge", "challenge")
	params.Set("code_challenge_method", "plain")

	// Construct the URL
	authURL := url.URL{
		Scheme:   "https",
		Host:     "twitter.com",
		Path:     "/i/oauth2/authorize",
		RawQuery: params.Encode(),
	}

	// build the authorization URL and open the browser for authorization
	openbrowser(authURL.String())

	// once authorized it'll redirect to localhost:3003/token need to handle that with a temporary webserver
	err := c.listenForToken()

	return err
}

func (c *Client) tokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/token" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	if !r.URL.Query().Has("code") {
		http.Error(w, "Received Token is empty", http.StatusNotAcceptable)
	}

	receivedToken := r.URL.Query().Get("code")

	c.oAuthToken = receivedToken

	logger.Info("Token received ✅")

	// Write success message to the response writer
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "<body style='text-align=center'><h1>RECEIVED TOKEN, YOU CAN CLOSE THIS WINDOW AND GO BACK TO THE CLI<h1></body>")

	go func() {
		time.Sleep(5 * time.Second)
		err := c.StopServer()
		if err != nil {
			fmt.Println(fmt.Errorf("Error stopping server: %v", err.Error()))
		}
	}()
}

func (c *Client) listenForToken() error {

	logger.Info("Listening for token...")

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) { c.tokenHandler(w, r) })

	//start a server at 3003
	err := c.ListenAndServe(":3003")

	return err
}

// Access Token generation

func (c *Client) generateBasicAuthToken() string {
	return fmt.Sprintf("Basic %v", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", c.clientId, c.clientSecret))))
}

func (c *Client) GenerateAccessToken() error {
	logger.Info("Generating access token")

	data := url.Values{}

	data.Set("code", c.oAuthToken)
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", c.clientId)
	data.Set("redirect_uri", "http://localhost:3003/token")
	data.Set("code_verifier", "challenge")

	r, err := http.NewRequest("POST", generateAccessTokenUrl, bytes.NewBuffer([]byte(data.Encode())))

	if err != nil {
		return fmt.Errorf("error creating request generateToken: %v", err)
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Authorization", c.generateBasicAuthToken())

	client := &http.Client{}
	resp, err := client.Do(r)

	if err != nil {
		return fmt.Errorf("error sending request generateToken: %v", err)
	}

	defer resp.Body.Close()

	accessTokenResponse := &AccessTokenResponse{}

	derr := json.NewDecoder(resp.Body).Decode(accessTokenResponse)

	if derr != nil {
		return fmt.Errorf("json decode error generateToken: %v", derr)
	}

	c.accessToken = accessTokenResponse

	logger.Info("Token generated, ready to tweet! 🚀")

	return nil
}
