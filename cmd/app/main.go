package main

import (
	"fmt"

	"github.com/nitishxyz/GoTweetCLI/internals/twitter"
)

// how to organise this project in one struct

func main() {

	client := twitter.NewClient("CLIENT_ID", "CLIENT_SECRET")

	client.AuthorizeAccount()

	errg := client.GenerateAccessToken()

	if errg != nil {
		fmt.Println(errg)
		return
	}

	err := client.PostTweet()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Tweet posted successfully!")
}
